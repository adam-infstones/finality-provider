package local

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/babylonchain/babylon/crypto/eots"
	bbntypes "github.com/babylonchain/babylon/types"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/go-bip39"

	"github.com/babylonchain/btc-validator/eotsmanager"
	"github.com/babylonchain/btc-validator/types"
	"github.com/babylonchain/btc-validator/valcfg"
)

const (
	secp256k1Type       = "secp256k1"
	mnemonicEntropySize = 256
)

type LocalEOTSManager struct {
	kr keyring.Keyring
	es *EOTSStore
}

func NewLocalEOTSManager(ctx client.Context, keyringBackend string, eotsCfg *valcfg.EOTSManagerConfig) (*LocalEOTSManager, error) {
	if keyringBackend == "" {
		return nil, fmt.Errorf("the keyring backend should not be empty")
	}

	kr, err := keyring.New(
		ctx.ChainID,
		// TODO currently only support test backend
		keyringBackend,
		ctx.KeyringDir,
		ctx.Input,
		ctx.Codec,
		ctx.KeyringOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create keyring: %w", err)
	}

	es, err := NewEOTSStore(eotsCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open the store for validators: %w", err)
	}

	return &LocalEOTSManager{
		kr: kr,
		es: es,
	}, nil
}

func (lm *LocalEOTSManager) CreateValidator(name, passPhrase string) ([]byte, error) {
	if lm.keyExists(name) {
		return nil, eotsmanager.ErrValidatorAlreadyExisted
	}

	keyringAlgos, _ := lm.kr.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(secp256k1Type, keyringAlgos)
	if err != nil {
		return nil, err
	}

	// read entropy seed straight from tmcrypto.Rand and convert to mnemonic
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return nil, err
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return nil, err
	}

	// TODO for now we leave and hdPath empty
	record, err := lm.kr.NewAccount(name, mnemonic, passPhrase, "", algo)
	if err != nil {
		return nil, err
	}

	pubKey, err := record.GetPubKey()
	if err != nil {
		return nil, err
	}

	var eotsPk *bbntypes.BIP340PubKey
	switch v := pubKey.(type) {
	case *secp256k1.PubKey:
		pk, err := btcec.ParsePubKey(v.Key)
		if err != nil {
			return nil, err
		}
		eotsPk = bbntypes.NewBIP340PubKeyFromBTCPK(pk)
	default:
		return nil, fmt.Errorf("unsupported key type in keyring")
	}

	if err := lm.es.saveValidatorKey(eotsPk.MustMarshal(), name); err != nil {
		return nil, err
	}

	return eotsPk.MustMarshal(), nil
}

func (lm *LocalEOTSManager) CreateRandomnessPairListWithExistenceCheck(valPk []byte, chainID []byte, startHeight uint64, num uint32) ([]*btcec.FieldVal, error) {
	// check whether the randomness is created already
	for i := uint32(0); i < num; i++ {
		height := startHeight + uint64(i)
		exists, err := lm.es.randPairExists(valPk, chainID, height)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, eotsmanager.ErrSchnorrRandomnessALreadyCreated
		}
	}

	return lm.CreateRandomnessPairList(valPk, chainID, startHeight, num)
}

func (lm *LocalEOTSManager) CreateRandomnessPairList(valPk []byte, chainID []byte, startHeight uint64, num uint32) ([]*btcec.FieldVal, error) {
	prList := make([]*btcec.FieldVal, 0, num)

	// TODO improve the security of randomness generation if concerned
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := uint32(0); i < num; i++ {
		// generate randomness pair
		eotsSR, eotsPR, err := eots.RandGen(r)
		if err != nil {
			return nil, fmt.Errorf("failed to generate randomness pair: %w", err)
		}

		// persists randomness pair
		height := startHeight + uint64(i)
		privRand := eotsSR.Bytes()
		pubRand := eotsPR.Bytes()
		randPair, err := types.NewSchnorrRandPair(privRand[:], pubRand[:])
		if err != nil {
			return nil, fmt.Errorf("invalid Schnorr randomness")
		}
		if err := lm.es.saveRandPair(valPk, chainID, height, randPair); err != nil {
			return nil, fmt.Errorf("failed to save randomness pair: %w", err)
		}

		prList = append(prList, eotsPR)
	}

	return prList, nil
}

func (lm *LocalEOTSManager) SignEOTS(valPk []byte, chainID []byte, msg []byte, height uint64) (*btcec.ModNScalar, error) {
	privRand, err := lm.getPrivRandomness(valPk, chainID, height)
	if err != nil {
		return nil, fmt.Errorf("failed to get private randomness: %w", err)
	}

	privKey, err := lm.getEOTSPrivKey(valPk)
	if err != nil {
		return nil, fmt.Errorf("failed to get EOTS private key: %w", err)
	}

	return eots.Sign(privKey, privRand, msg)
}

func (lm *LocalEOTSManager) SignSchnorrSig(valPk []byte, msg []byte) (*schnorr.Signature, error) {
	privKey, err := lm.getEOTSPrivKey(valPk)
	if err != nil {
		return nil, fmt.Errorf("failed to get EOTS private key: %w", err)
	}

	return schnorr.Sign(privKey, msg)
}

func (lm *LocalEOTSManager) Close() error {
	return lm.es.Close()
}

func (lm *LocalEOTSManager) getPrivRandomness(valPk []byte, chainID []byte, height uint64) (*eots.PrivateRand, error) {
	randPair, err := lm.es.getRandPair(valPk, chainID, height)
	if err != nil {
		return nil, err
	}

	if len(randPair.SecRand) != types.SchnorrRandomnessLength {
		return nil, fmt.Errorf("the private randomness should be 32 bytes")
	}

	privRand := new(eots.PrivateRand)
	privRand.SetByteSlice(randPair.SecRand)

	return privRand, nil
}

func (lm *LocalEOTSManager) GetValidatorKeyName(valPk []byte) (string, error) {
	return lm.es.getValidatorKeyName(valPk)
}

func (lm *LocalEOTSManager) getEOTSPrivKey(valPk []byte) (*btcec.PrivateKey, error) {
	keyName, err := lm.es.getValidatorKeyName(valPk)
	if err != nil {
		return nil, err
	}

	k, err := lm.kr.Key(keyName)
	if err != nil {
		return nil, err
	}

	privKeyCached := k.GetLocal().PrivKey.GetCachedValue()

	var privKey *btcec.PrivateKey
	switch v := privKeyCached.(type) {
	case *secp256k1.PrivKey:
		privKey, _ = btcec.PrivKeyFromBytes(v.Key)
		return privKey, nil
	default:
		return nil, fmt.Errorf("unsupported key type in keyring")
	}
}

func (lm *LocalEOTSManager) keyExists(name string) bool {
	_, err := lm.kr.Key(name)
	return err == nil
}
package service_test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/babylonchain/babylon/testutil/datagen"
	bstypes "github.com/babylonchain/babylon/x/btcstaking/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	cometbfttypes "github.com/cometbft/cometbft/types"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/babylonchain/btc-validator/clientcontroller"
	"github.com/babylonchain/btc-validator/proto"
	"github.com/babylonchain/btc-validator/service"
	"github.com/babylonchain/btc-validator/testutil"
	"github.com/babylonchain/btc-validator/testutil/mocks"
	"github.com/babylonchain/btc-validator/val"
	"github.com/babylonchain/btc-validator/valcfg"
)

var (
	eventuallyWaitTimeOut = 1 * time.Second
	eventuallyPollTime    = 10 * time.Millisecond
)

func FuzzStatusUpdate(f *testing.F) {
	testutil.AddRandomSeedsToFuzzer(f, 10)
	f.Fuzz(func(t *testing.T, seed int64) {
		r := rand.New(rand.NewSource(seed))

		ctl := gomock.NewController(t)
		mockClientController := mocks.NewMockClientController(ctl)
		vm, cleanUp := newValidatorManagerWithRegisteredValidator(t, r, mockClientController)
		defer cleanUp()

		// setup mocks
		currentHeight := uint64(r.Int63n(100) + 1)
		currentHeaderRes := &coretypes.ResultHeader{
			Header: &cometbfttypes.Header{
				Height:         int64(currentHeight),
				LastCommitHash: datagen.GenRandomByteArray(r, 32),
			},
		}
		mockClientController.EXPECT().QueryBestHeader().Return(currentHeaderRes, nil).AnyTimes()
		status := &coretypes.ResultStatus{
			SyncInfo: coretypes.SyncInfo{LatestBlockHeight: int64(currentHeight)},
		}
		mockClientController.EXPECT().QueryNodeStatus().Return(status, nil).AnyTimes()
		mockClientController.EXPECT().Close().Return(nil).AnyTimes()
		mockClientController.EXPECT().QueryLatestFinalizedBlocks(gomock.Any()).Return(nil, nil).AnyTimes()
		mockClientController.EXPECT().QueryBestHeader().Return(currentHeaderRes, nil).AnyTimes()
		mockClientController.EXPECT().QueryHeader(gomock.Any()).Return(currentHeaderRes, nil).AnyTimes()

		votingPower := uint64(r.Intn(2))
		mockClientController.EXPECT().QueryValidatorVotingPower(gomock.Any(), currentHeight).Return(votingPower, nil).AnyTimes()
		var slashedHeight uint64
		if votingPower == 0 {
			slashedHeight = uint64(r.Intn(2))
			btcVal := &bstypes.BTCValidator{SlashedBtcHeight: slashedHeight}
			mockClientController.EXPECT().QueryValidator(gomock.Any()).Return(btcVal, nil).AnyTimes()
		}

		err := vm.Start()
		require.NoError(t, err)
		valIns := vm.ListValidatorInstances()[0]
		// stop the validator as we are testing static functionalities
		err = valIns.Stop()
		require.NoError(t, err)

		if votingPower > 0 {
			waitForStatus(t, valIns, proto.ValidatorStatus_ACTIVE)
		} else {
			if slashedHeight == 0 && valIns.GetStatus() == proto.ValidatorStatus_ACTIVE {
				waitForStatus(t, valIns, proto.ValidatorStatus_INACTIVE)
			} else if slashedHeight > 0 {
				waitForStatus(t, valIns, proto.ValidatorStatus_SLASHED)
			}
		}
	})
}

func waitForStatus(t *testing.T, valIns *service.ValidatorInstance, s proto.ValidatorStatus) {
	require.Eventually(t,
		func() bool {
			return valIns.GetStatus() == s
		}, eventuallyWaitTimeOut, eventuallyPollTime)
}

func newValidatorManagerWithRegisteredValidator(t *testing.T, r *rand.Rand, cc clientcontroller.ClientController) (*service.ValidatorManager, func()) {
	// create validator app with config
	cfg := valcfg.DefaultConfig()
	cfg.StatusUpdateInterval = 10 * time.Millisecond
	cfg.DatabaseConfig = testutil.GenDBConfig(r, t)
	cfg.BabylonConfig.KeyDirectory = t.TempDir()
	logger := logrus.New()

	kr, err := service.CreateKeyring(cfg.BabylonConfig.KeyDirectory,
		cfg.BabylonConfig.ChainID,
		cfg.BabylonConfig.KeyringBackend)
	require.NoError(t, err)

	valStore, err := val.NewValidatorStore(cfg.DatabaseConfig)
	require.NoError(t, err)

	vm, err := service.NewValidatorManager(valStore, &cfg, kr, cc, logger)
	require.NoError(t, err)

	// create registered validator
	keyName := datagen.GenRandomHexStr(r, 10)
	kc, err := val.NewKeyringControllerWithKeyring(kr, keyName)
	require.NoError(t, err)
	btcPk, bbnPk, err := kc.CreateValidatorKeys()
	require.NoError(t, err)
	pop, err := kc.CreatePop()
	require.NoError(t, err)

	storedValidator := val.NewStoreValidator(bbnPk, btcPk, keyName, pop, testutil.EmptyDescription(), testutil.ZeroCommissionRate())
	storedValidator.Status = proto.ValidatorStatus_REGISTERED
	err = valStore.SaveValidator(storedValidator)
	require.NoError(t, err)

	cleanUp := func() {
		err = vm.Stop()
		require.NoError(t, err)
		err := os.RemoveAll(cfg.DatabaseConfig.Path)
		require.NoError(t, err)
		err = os.RemoveAll(cfg.BabylonConfig.KeyDirectory)
		require.NoError(t, err)
	}

	return vm, cleanUp
}
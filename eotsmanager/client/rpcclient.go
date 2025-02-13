package client

import (
	"context"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/babylonchain/finality-provider/eotsmanager"
	"github.com/babylonchain/finality-provider/eotsmanager/proto"
	"github.com/babylonchain/finality-provider/eotsmanager/types"
)

var _ eotsmanager.EOTSManager = &EOTSManagerGRpcClient{}

type EOTSManagerGRpcClient struct {
	client proto.EOTSManagerClient
	conn   *grpc.ClientConn
}

func NewEOTSManagerGRpcClient(remoteAddr string) (*EOTSManagerGRpcClient, error) {
	conn, err := grpc.Dial(remoteAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to build gRPC connection to %s: %w", remoteAddr, err)
	}

	gClient := &EOTSManagerGRpcClient{
		client: proto.NewEOTSManagerClient(conn),
		conn:   conn,
	}

	if err := gClient.Ping(); err != nil {
		return nil, fmt.Errorf("the EOTS manager server is not responding: %w", err)
	}

	return gClient, nil
}

func (c *EOTSManagerGRpcClient) Ping() error {
	req := &proto.PingRequest{}

	_, err := c.client.Ping(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}

func (c *EOTSManagerGRpcClient) CreateKey(name, passphrase, hdPath string) ([]byte, error) {
	req := &proto.CreateKeyRequest{Name: name, Passphrase: passphrase, HdPath: hdPath}
	res, err := c.client.CreateKey(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return res.Pk, nil
}

func (c *EOTSManagerGRpcClient) CreateMasterRandPair(uid, chainID []byte, passphrase string) (string, error) {
	req := &proto.CreateMasterRandPairRequest{
		Uid:        uid,
		ChainId:    chainID,
		Passphrase: passphrase,
	}
	res, err := c.client.CreateMasterRandPair(context.Background(), req)
	if err != nil {
		return "", err
	}

	return res.MasterPubRand, nil
}

func (c *EOTSManagerGRpcClient) KeyRecord(uid []byte, passphrase string) (*types.KeyRecord, error) {
	req := &proto.KeyRecordRequest{Uid: uid, Passphrase: passphrase}

	res, err := c.client.KeyRecord(context.Background(), req)
	if err != nil {
		return nil, err
	}

	privKey, _ := btcec.PrivKeyFromBytes(res.PrivateKey)

	return &types.KeyRecord{
		Name:    res.Name,
		PrivKey: privKey,
	}, nil
}

func (c *EOTSManagerGRpcClient) SignEOTS(uid, chaiID, msg []byte, height uint64, passphrase string) (*btcec.ModNScalar, error) {
	req := &proto.SignEOTSRequest{
		Uid:        uid,
		ChainId:    chaiID,
		Msg:        msg,
		Height:     height,
		Passphrase: passphrase,
	}
	res, err := c.client.SignEOTS(context.Background(), req)
	if err != nil {
		return nil, err
	}

	var s btcec.ModNScalar
	s.SetByteSlice(res.Sig)

	return &s, nil
}

func (c *EOTSManagerGRpcClient) SignSchnorrSig(uid, msg []byte, passphrase string) (*schnorr.Signature, error) {
	req := &proto.SignSchnorrSigRequest{Uid: uid, Msg: msg, Passphrase: passphrase}
	res, err := c.client.SignSchnorrSig(context.Background(), req)
	if err != nil {
		return nil, err
	}

	sig, err := schnorr.ParseSignature(res.Sig)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

func (c *EOTSManagerGRpcClient) Close() error {
	return c.conn.Close()
}

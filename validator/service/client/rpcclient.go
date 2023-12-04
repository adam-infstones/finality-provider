package client

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	bbntypes "github.com/babylonchain/babylon/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/babylonchain/btc-validator/validator/proto"
)

type ValidatorServiceGRpcClient struct {
	client proto.BtcValidatorsClient
}

func NewValidatorServiceGRpcClient(remoteAddr string) (*ValidatorServiceGRpcClient, func(), error) {
	conn, err := grpc.Dial(remoteAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build gRPC connection to %s: %w", remoteAddr, err)
	}

	cleanUp := func() {
		conn.Close()
	}

	return &ValidatorServiceGRpcClient{
		client: proto.NewBtcValidatorsClient(conn),
	}, cleanUp, nil
}

func (c *ValidatorServiceGRpcClient) GetInfo(ctx context.Context) (*proto.GetInfoResponse, error) {
	req := &proto.GetInfoRequest{}
	res, err := c.client.GetInfo(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ValidatorServiceGRpcClient) RegisterValidator(
	ctx context.Context,
	valPk *bbntypes.BIP340PubKey,
	passphrase string,
) (*proto.RegisterValidatorResponse, error) {

	req := &proto.RegisterValidatorRequest{BtcPk: valPk.MarshalHex(), Passphrase: passphrase}
	res, err := c.client.RegisterValidator(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ValidatorServiceGRpcClient) CreateValidator(
	ctx context.Context,
	keyName, chainID, passphrase, hdPath string,
	description types.Description,
	commission *sdkmath.LegacyDec,
) (*proto.CreateValidatorResponse, error) {

	descBytes, err := description.Marshal()
	if err != nil {
		return nil, err
	}

	req := &proto.CreateValidatorRequest{
		KeyName:     keyName,
		ChainId:     chainID,
		Passphrase:  passphrase,
		HdPath:      hdPath,
		Description: descBytes,
		Commission:  commission.String(),
	}

	res, err := c.client.CreateValidator(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ValidatorServiceGRpcClient) AddFinalitySignature(ctx context.Context, valPk string, height uint64, lch []byte) (*proto.AddFinalitySignatureResponse, error) {
	req := &proto.AddFinalitySignatureRequest{
		BtcPk:   valPk,
		Height:  height,
		AppHash: lch,
	}

	res, err := c.client.AddFinalitySignature(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ValidatorServiceGRpcClient) QueryValidatorList(ctx context.Context) (*proto.QueryValidatorListResponse, error) {
	req := &proto.QueryValidatorListRequest{}
	res, err := c.client.QueryValidatorList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ValidatorServiceGRpcClient) QueryValidatorInfo(ctx context.Context, valPk *bbntypes.BIP340PubKey) (*proto.QueryValidatorResponse, error) {
	req := &proto.QueryValidatorRequest{BtcPk: valPk.MarshalHex()}
	res, err := c.client.QueryValidator(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
// Code generated by MockGen. DO NOT EDIT.
// Source: bbnclient/interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	finalitytypes "github.com/babylonchain/babylon/x/finality/types"
	reflect "reflect"

	types "github.com/babylonchain/babylon/types"
	types0 "github.com/babylonchain/babylon/x/btcstaking/types"
	babylonclient "github.com/babylonchain/btc-validator/bbnclient"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	gomock "github.com/golang/mock/gomock"
)

// MockBabylonClient is a mock of BabylonClient interface.
type MockBabylonClient struct {
	ctrl     *gomock.Controller
	recorder *MockBabylonClientMockRecorder
}

// MockBabylonClientMockRecorder is the mock recorder for MockBabylonClient.
type MockBabylonClientMockRecorder struct {
	mock *MockBabylonClient
}

// NewMockBabylonClient creates a new mock instance.
func NewMockBabylonClient(ctrl *gomock.Controller) *MockBabylonClient {
	mock := &MockBabylonClient{ctrl: ctrl}
	mock.recorder = &MockBabylonClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBabylonClient) EXPECT() *MockBabylonClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockBabylonClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockBabylonClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockBabylonClient)(nil).Close))
}

// CommitPubRandList mocks base method.
func (m *MockBabylonClient) CommitPubRandList(btcPubKey *types.BIP340PubKey, startHeight uint64, pubRandList []types.SchnorrPubRand, sig *types.BIP340Signature) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitPubRandList", btcPubKey, startHeight, pubRandList, sig)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CommitPubRandList indicates an expected call of CommitPubRandList.
func (mr *MockBabylonClientMockRecorder) CommitPubRandList(btcPubKey, startHeight, pubRandList, sig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitPubRandList", reflect.TypeOf((*MockBabylonClient)(nil).CommitPubRandList), btcPubKey, startHeight, pubRandList, sig)
}

// GetStakingParams mocks base method.
func (m *MockBabylonClient) GetStakingParams() (*babylonclient.StakingParams, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStakingParams")
	ret0, _ := ret[0].(*babylonclient.StakingParams)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStakingParams indicates an expected call of GetStakingParams.
func (mr *MockBabylonClientMockRecorder) GetStakingParams() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStakingParams", reflect.TypeOf((*MockBabylonClient)(nil).GetStakingParams))
}

// QueryBestHeader mocks base method.
func (m *MockBabylonClient) QueryBestHeader() (*coretypes.ResultHeader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryBestHeader")
	ret0, _ := ret[0].(*coretypes.ResultHeader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryBestHeader indicates an expected call of QueryBestHeader.
func (mr *MockBabylonClientMockRecorder) QueryBestHeader() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryBestHeader", reflect.TypeOf((*MockBabylonClient)(nil).QueryBestHeader))
}

// QueryHeader mocks base method.
func (m *MockBabylonClient) QueryHeader(height int64) (*coretypes.ResultHeader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryHeader", height)
	ret0, _ := ret[0].(*coretypes.ResultHeader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryHeader indicates an expected call of QueryHeader.
func (mr *MockBabylonClientMockRecorder) QueryHeader(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryHeader", reflect.TypeOf((*MockBabylonClient)(nil).QueryHeader), height)
}

// QueryHeightWithLastPubRand mocks base method.
func (m *MockBabylonClient) QueryHeightWithLastPubRand(btcPubKey *types.BIP340PubKey) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryHeightWithLastPubRand", btcPubKey)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryHeightWithLastPubRand indicates an expected call of QueryHeightWithLastPubRand.
func (mr *MockBabylonClientMockRecorder) QueryHeightWithLastPubRand(btcPubKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryHeightWithLastPubRand", reflect.TypeOf((*MockBabylonClient)(nil).QueryHeightWithLastPubRand), btcPubKey)
}

// QueryLatestFinalisedBlocks mocks base method.
func (m *MockBabylonClient) QueryLatestFinalisedBlocks(count uint64) ([]*finalitytypes.IndexedBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryLatestFinalisedBlocks", count)
	ret0, _ := ret[0].([]*finalitytypes.IndexedBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryLatestFinalisedBlocks indicates an expected call of QueryLatestFinalisedBlocks.
func (mr *MockBabylonClientMockRecorder) QueryLatestFinalisedBlocks(count interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryLatestFinalisedBlocks", reflect.TypeOf((*MockBabylonClient)(nil).QueryLatestFinalisedBlocks), count)
}

// QueryNodeStatus mocks base method.
func (m *MockBabylonClient) QueryNodeStatus() (*coretypes.ResultStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryNodeStatus")
	ret0, _ := ret[0].(*coretypes.ResultStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryNodeStatus indicates an expected call of QueryNodeStatus.
func (mr *MockBabylonClientMockRecorder) QueryNodeStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryNodeStatus", reflect.TypeOf((*MockBabylonClient)(nil).QueryNodeStatus))
}

// QueryPendingBTCDelegations mocks base method.
func (m *MockBabylonClient) QueryPendingBTCDelegations() ([]*types0.BTCDelegation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryPendingBTCDelegations")
	ret0, _ := ret[0].([]*types0.BTCDelegation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryPendingBTCDelegations indicates an expected call of QueryPendingBTCDelegations.
func (mr *MockBabylonClientMockRecorder) QueryPendingBTCDelegations() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryPendingBTCDelegations", reflect.TypeOf((*MockBabylonClient)(nil).QueryPendingBTCDelegations))
}

// QueryValidatorVotingPower mocks base method.
func (m *MockBabylonClient) QueryValidatorVotingPower(btcPubKey *types.BIP340PubKey, blockHeight uint64) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryValidatorVotingPower", btcPubKey, blockHeight)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryValidatorVotingPower indicates an expected call of QueryValidatorVotingPower.
func (mr *MockBabylonClientMockRecorder) QueryValidatorVotingPower(btcPubKey, blockHeight interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryValidatorVotingPower", reflect.TypeOf((*MockBabylonClient)(nil).QueryValidatorVotingPower), btcPubKey, blockHeight)
}

// RegisterValidator mocks base method.
func (m *MockBabylonClient) RegisterValidator(bbnPubKey *secp256k1.PubKey, btcPubKey *types.BIP340PubKey, pop *types0.ProofOfPossession) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterValidator", bbnPubKey, btcPubKey, pop)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterValidator indicates an expected call of RegisterValidator.
func (mr *MockBabylonClientMockRecorder) RegisterValidator(bbnPubKey, btcPubKey, pop interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterValidator", reflect.TypeOf((*MockBabylonClient)(nil).RegisterValidator), bbnPubKey, btcPubKey, pop)
}

// SubmitFinalitySig mocks base method.
func (m *MockBabylonClient) SubmitFinalitySig(btcPubKey *types.BIP340PubKey, blockHeight uint64, blockHash []byte, sig *types.SchnorrEOTSSig) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitFinalitySig", btcPubKey, blockHeight, blockHash, sig)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitFinalitySig indicates an expected call of SubmitFinalitySig.
func (mr *MockBabylonClientMockRecorder) SubmitFinalitySig(btcPubKey, blockHeight, blockHash, sig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitFinalitySig", reflect.TypeOf((*MockBabylonClient)(nil).SubmitFinalitySig), btcPubKey, blockHeight, blockHash, sig)
}

// SubmitJurySig mocks base method.
func (m *MockBabylonClient) SubmitJurySig(btcPubKey, delPubKey *types.BIP340PubKey, sig *types.BIP340Signature) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitJurySig", btcPubKey, delPubKey, sig)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitJurySig indicates an expected call of SubmitJurySig.
func (mr *MockBabylonClientMockRecorder) SubmitJurySig(btcPubKey, delPubKey, sig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitJurySig", reflect.TypeOf((*MockBabylonClient)(nil).SubmitJurySig), btcPubKey, delPubKey, sig)
}

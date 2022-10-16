package rpc

import (
	"WeCross-Go-SDK/logger"
	"WeCross-Go-SDK/rpc/eles/account"
	"WeCross-Go-SDK/rpc/service"
	"WeCross-Go-SDK/rpc/types"
	"WeCross-Go-SDK/rpc/types/response"
)

const WeCrossRPCModelTag = "WeCrossRPCModel"

type WeCrossRPCModel struct {
	logger         *logger.Logger
	weCrossService service.WeCrossService
}

func NewWeCrossRPCModel(weCrossService service.WeCrossService) *WeCrossRPCModel {
	logger := logger.NewLogger(WeCrossRPCModelTag)
	return &WeCrossRPCModel{
		logger:         logger,
		weCrossService: weCrossService,
	}
}

func (w WeCrossRPCModel) Test() *RemoteCall {
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            "/sys/test",
		responseType:   response.CommonResponse,
		request:        types.NewRequest(nil),
	}
}

func (w WeCrossRPCModel) SupportedStubs() *RemoteCall {
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "GET",
		uri:            "/sys/supportedStubs",
		responseType:   response.StubResponse,
		request:        types.NewRequest(nil),
	}
}

func (w WeCrossRPCModel) QueryPub() *RemoteCall {
	return &RemoteCall{}
}

func (w WeCrossRPCModel) QueryAuthCode() *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) ListAccount() *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) ListResources(ignoreRemote bool) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) Detail(path string) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) Call(path, method string, args ...string) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) SendTransaction(path, method string, args ...string) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) Invoke(path, method string, args ...string) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) CallXA(transactionID, path, method string, args ...string) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) SendXATransaction(transactionID, path, method string, args ...string) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) StartXATransaction(transactionID string, paths []string) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) CommitXATransaction(transactionID string, paths []string) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) RollbackXATransaction(transactionID string, paths []string) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) GetXATransaction(transactionID string, paths []string) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) CustomCommand(command string, path string, args ...any) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) ListXATransactions(size int) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) Register(name, password string) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) Login(name, password string) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) Logout() *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) AddChainAccount(chainType string, chainAccount *account.ChainAccount) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) SetDefaultAccount(chainType string, chainAccount *account.ChainAccount, keyID int) *RemoteCall {
	//TODO implement me
	panic("implement me")
}

func (w WeCrossRPCModel) GetCurrentTransactionID() string {
	//TODO implement me
	panic("implement me")
}

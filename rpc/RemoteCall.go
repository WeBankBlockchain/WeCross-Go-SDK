package rpc

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/common"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/service"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/types"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/types/response"
)

type RemoteCall struct {
	weCrossService service.WeCrossService
	httpMethod     string
	uri            string
	responseType   response.ResponseType
	request        *types.Request
}

func NewRemoteCall(weCrossService service.WeCrossService, httpMethod string, uri string, responseType response.ResponseType, request *types.Request) *RemoteCall {
	return &RemoteCall{
		weCrossService: weCrossService,
		httpMethod:     httpMethod,
		uri:            uri,
		responseType:   responseType,
		request:        request,
	}
}

func (rc *RemoteCall) Send() (*types.Response, *common.WeCrossSDKError) {
	return rc.weCrossService.Send(rc.httpMethod, rc.uri, rc.request, rc.responseType)
}

func (rc *RemoteCall) AsyncSend(back *types.CallBack) {
	rc.weCrossService.AsyncSend(rc.httpMethod, rc.uri, rc.request, rc.responseType, back)
}

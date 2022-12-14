package service

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/common"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/service/transactionContext"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/types"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/types/response"
)

type WeCrossService interface {
	Init() *common.WeCrossSDKError
	Send(httpMethod string, uri string, request *types.Request, responseType response.ResponseType) (*types.Response, *common.WeCrossSDKError)
	AsyncSend(httpMethod string, uri string, request *types.Request, responseType response.ResponseType, back *types.CallBack)
	GetTransactionContex() *transactionContext.TxCtx
}

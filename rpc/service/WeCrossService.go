package service

import (
	"WeCross-Go-SDK/common"
	"WeCross-Go-SDK/rpc/types"
	"WeCross-Go-SDK/rpc/types/response"
)

type WeCrossService interface {
	Init() *common.WeCrossSDKError
	Send(httpMethod string, uri string, request *types.Request, responseType response.ResponseType) *types.Response
	AsynSend(httpMethod string, uri string, request *types.Request, responseType response.ResponseType) *types.Response
}

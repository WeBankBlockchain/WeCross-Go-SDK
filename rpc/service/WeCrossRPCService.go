package service

import (
	"WeCross-Go-SDK/common"
	"WeCross-Go-SDK/logger"
	"WeCross-Go-SDK/rpc/service/httpAsyncClient"
	"WeCross-Go-SDK/rpc/types"
	"WeCross-Go-SDK/rpc/types/response"
	"WeCross-Go-SDK/utils"
)

type WeCrossRPCService struct {
	logger *logger.Logger

	server     string
	httpClient *httpAsyncClient.AsyncHttpClient
	urlPrefix  string
}

func (wcs *WeCrossRPCService) Init() *common.WeCrossSDKError {
	config, err := utils.GetToml(common.APPLICATION_CONFIG_FILE)
	if err != nil {
		return err
	}
	connection, err := NewConnection(config)
	if err != nil {
		return err
	}
	wcs.logger.Infof("RPCService init: %s", connection.ToString())

	// set server
	server := ""
	if connection.GetSslSwitch() == common.SSL_OFF {
		server += "http://"
	} else {
		server += "https://"
	}
	server += connection.GetServer()
	wcs.server = server

	// set urlPrefix
	if len(connection.GetUrlPrefix()) != 0 {
		wcs.urlPrefix = connection.GetUrlPrefix()
	}

	// set httpClient
	httpClient, err := httpAsyncClient.NewAsyncHttpClient(connection)
	if err != nil {
		return err
	}
	wcs.httpClient = httpClient

	return nil
}

func (wcs *WeCrossRPCService) Send(httpMethod string, uri string, request *types.Request, responseType response.ResponseType) *types.Response {
	//TODO implement me
	panic("implement me")
}

func (wcs *WeCrossRPCService) AsynSend(httpMethod string, uri string, request *types.Request, responseType response.ResponseType) *types.Response {
	//TODO implement me
	panic("implement me")
}

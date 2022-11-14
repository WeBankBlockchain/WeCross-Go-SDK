package service

import (
	"WeCross-Go-SDK/common"
	"WeCross-Go-SDK/logger"
	"WeCross-Go-SDK/rpc/eles"
	"WeCross-Go-SDK/rpc/service/authentication"
	"WeCross-Go-SDK/rpc/service/connections"
	"WeCross-Go-SDK/rpc/service/httpAsyncClient"
	"WeCross-Go-SDK/rpc/service/transactionContext"
	"WeCross-Go-SDK/rpc/types"
	"WeCross-Go-SDK/rpc/types/request"
	"WeCross-Go-SDK/rpc/types/response"
	"WeCross-Go-SDK/utils"
	"fmt"
	"strings"
	"time"
)

const (
	WeCrossRPCServiceTag = "WeCrossRPCService"

	Max_Send_Wait_Time = 20 * time.Second
)

type WeCrossRPCService struct {
	logger *logger.Logger

	server     string
	httpClient *httpAsyncClient.AsyncHttpClient
	urlPrefix  string

	authenticationManager *authentication.AuthenticationManager
	transactionContext    *transactionContext.TransactionContext
}

func NewWeCrossRPCService() *WeCrossRPCService {
	logger := logger.NewLogger(WeCrossRPCServiceTag)
	return &WeCrossRPCService{
		logger:                logger,
		authenticationManager: authentication.NewAuthenticationManager(),
		transactionContext:    transactionContext.GlobalTransactionContext,
	}
}

func (wcs *WeCrossRPCService) Init() *common.WeCrossSDKError {
	config, err := utils.GetToml(utils.ReadClassPath(common.APPLICATION_CONFIG_FILE))
	if err != nil {
		return err
	}
	connection, err := connections.NewConnection(config)
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

func (wcs *WeCrossRPCService) Send(httpMethod string, uri string, inputRequest *types.Request, responseType response.ResponseType) (*types.Response, *common.WeCrossSDKError) {
	var finalResponse *types.Response
	var finalError *common.WeCrossSDKError
	finish := make(chan struct{})
	onSuccess := func(inResponse *types.Response) {
		finalResponse = inResponse
		select {
		case finish <- struct{}{}:
			return
		default:
			return
		}
	}
	onFailed := func(err *common.WeCrossSDKError) {
		wcs.logger.Warnf("send onFailed: %s", err.Content)
		finalError = err
		select {
		case finish <- struct{}{}:
			return
		default:
			return
		}
	}

	callBack := types.NewCallBack(onSuccess, onFailed)
	wcs.AsyncSend(httpMethod, uri, inputRequest, responseType, callBack)
	outTimer := time.NewTimer(Max_Send_Wait_Time)
	select {
	case <-finish:
		wcs.logger.Debugf("response: %s", finalResponse.ToString())
		if finalError != nil {
			return finalResponse, finalError
		}

		if responseType == response.UAResponse {
			if uaRq, ok := inputRequest.Data.(*request.UARequest); ok {
				err := wcs.getUAResponseInfo(uri, uaRq, finalResponse)
				if err != nil {
					return finalResponse, err
				}
			} else if uaRq, ok := inputRequest.Ext.(*request.UARequest); ok {
				err := wcs.getUAResponseInfo(uri, uaRq, finalResponse)
				if err != nil {
					return finalResponse, err
				}
			}
		}

		if responseType == response.XAResponse {
			wcs.getXAResponseInfo(uri, inputRequest, finalResponse)
		}

		return finalResponse, finalError
	case <-outTimer.C:
		return nil, common.NewWeCrossSDKFromString(common.RPC_ERROR, "fain in Send: time out while waiting for response")
	}
}

func (wcs *WeCrossRPCService) AsyncSend(httpMethod string, uri string, request *types.Request, responseType response.ResponseType, back *types.CallBack) {
	url := wcs.server + wcs.urlPrefix + uri
	wcs.logger.Debugf("request: %s; url: %s", request.ToJson(), url)
	err := wcs.checkRequest(request)
	if err != nil {
		wcs.logger.Errorf(err.Error())
		back.CallOnFailed(common.NewWeCrossSDKFromString(common.INTERNAL_ERROR, "fail in AsyncSend:"+err.Content))
		return
	}

	httpRequest, err := wcs.httpClient.Prepare(strings.ToUpper(httpMethod), url, request.ToJson())
	if err != nil {
		wcs.logger.Errorf(err.Error())
		back.CallOnFailed(common.NewWeCrossSDKFromString(common.INTERNAL_ERROR, "fail in AsyncSend:"+err.Content))
		return
	}

	currentUserCredential := wcs.authenticationManager.GetCurrentUserCredential()

	method := utils.GetUriMethod(uri)
	wcs.logger.Debugf("uri path: %s", method)

	if _, ok := eles.AUTH_REQUIRED_COMMANDS[method]; ok {
		if len(currentUserCredential) == 0 {
			wcs.logger.Errorf("Request's method required AUTH, but current credential is null.")
			back.CallOnFailed(common.NewWeCrossSDKFromString(common.LACK_AUTHENTICATION, "Command "+method+" needs Auth, please login."))
			return
		}
		httpRequest.Header.Set("Authorization", currentUserCredential)
	}

	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")

	go func() {
		httpResponse, err := wcs.httpClient.SendRequest(httpRequest)
		if err != nil {
			back.CallOnFailed(common.NewWeCrossSDKFromString(common.RPC_ERROR, "fail in SendRequest:"+err.Content))
			return
		}
		if httpResponse.StatusCode == 401 {
			errContent := "HTTP status code: 401-Unauthorized, have you logged in?\n" + "If you have logged-in already, maybe you should re-login " + "because your account login status has expired."
			back.CallOnFailed(common.NewWeCrossSDKFromString(common.LACK_AUTHENTICATION, errContent))
			return
		}
		if httpResponse.StatusCode == 404 {
			errContent := "HTTP status code: 404 Not Found\n" + "Maybe your request's resource path is wrong."
			back.CallOnFailed(common.NewWeCrossSDKFromString(common.LACK_AUTHENTICATION, errContent))
			return
		}
		if httpResponse.StatusCode != 200 {
			errContent := fmt.Sprintf("HTTP response status: %d message: %s", httpResponse.StatusCode, httpResponse.Status)
			back.CallOnFailed(common.NewWeCrossSDKFromString(common.RPC_ERROR, errContent))
			return
		} else {
			gotResponse := types.ParseResponse(httpResponse, responseType)
			back.CallOnSuccess(gotResponse)
			return
		}
	}()

}

func (wcs *WeCrossRPCService) checkRequest(request *types.Request) *common.WeCrossSDKError {
	if request == nil {
		return common.NewWeCrossSDKFromString(common.RPC_ERROR, "Request is nil")
	}
	if len(request.GetVersion()) == 0 {
		return common.NewWeCrossSDKFromString(common.RPC_ERROR, "Request version is empty")
	}
	return nil
}

func (wcs *WeCrossRPCService) getUAResponseInfo(uri string, uaRequest *request.UARequest, uaResponse *types.Response) *common.WeCrossSDKError {
	query := utils.GetUriQuery(uri)
	if query == "login" {
		credential := ""
		uaRp, ok := uaResponse.Data.(*response.UAReceipt)
		if ok {
			credential = uaRp.Credential
		}
		wcs.logger.Infof("CurrentUser: %s", uaRequest.UserName)
		if len(credential) == 0 {
			wcs.logger.Errorln("Login fail, credential in UAResponse is null")
			return common.NewWeCrossSDKFromString(common.RPC_ERROR, "Login fail, credential in UAResponse is null!")
		}
		wcs.authenticationManager.SetCurrentUser(uaRequest.UserName, credential)
	}
	if query == "logout" {
		wcs.logger.Infof("CurrentUser: %s logout.", wcs.authenticationManager.GetCurrentUser())
		wcs.authenticationManager.ClearCurrentUser()
	}
	return nil
}

func (wcs *WeCrossRPCService) getXAResponseInfo(uri string, inRequest *types.Request, xaResponse *types.Response) {
	query := utils.GetUriQuery(uri)
	if query == "startXATransaction" && xaResponse.ErrorCode == 0 {
		xaRequest, ok := inRequest.Data.(*request.XATransactionRequest)
		if !ok {
			return
		}
		txCtx := wcs.transactionContext.GetContex()
		txCtx.TxID = xaRequest.XaTransactionID
		txCtx.Seq = 1
		txCtx.PathInTransaction = xaRequest.Paths
		wcs.transactionContext.SetContex(txCtx)
	} else if (query == "commitXATransaction") || query == "rollbackXATransaction" && xaResponse.ErrorCode == 0 {
		wcs.transactionContext.ClearAll()
	}
}

package rpc

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/common"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/common/cryptos"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/logger"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/eles"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/eles/account"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/service"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/types"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/types/request"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/types/response"
	"regexp"
	"strings"
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

func (w *WeCrossRPCModel) Test() *RemoteCall {
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            "/sys/test",
		responseType:   response.CommonResponse,
		request:        types.NewRequest(nil),
	}
}

func (w *WeCrossRPCModel) SupportedStubs() *RemoteCall {
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "GET",
		uri:            "/sys/supportedStubs",
		responseType:   response.StubResponse,
		request:        types.NewRequest(nil),
	}
}

func (w *WeCrossRPCModel) QueryPub() *RemoteCall {
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "GET",
		uri:            "/auth/pub",
		responseType:   response.PubResponse,
		request:        types.NewRequest(nil),
	}
}

func (w *WeCrossRPCModel) QueryAuthCode() *RemoteCall {
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "GET",
		uri:            "/auth/authCode",
		responseType:   response.AuthCodeResponse,
		request:        types.NewRequest(nil),
	}
}

func (w *WeCrossRPCModel) ListAccount() *RemoteCall {
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            "/auth/listAccount",
		responseType:   response.AccountResponse,
		request:        types.NewRequest(nil),
	}
}

func (w *WeCrossRPCModel) ListResources(ignoreRemote bool) *RemoteCall {
	resourceRequest := request.NewResourceRequest(ignoreRemote)
	wrappedRequest := types.NewRequest(resourceRequest)
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "GET",
		uri:            "/sys/listResources",
		responseType:   response.ResourceResponse,
		request:        wrappedRequest,
	}
}

func (w *WeCrossRPCModel) Detail(path string) *RemoteCall {
	uri := "/resource/" + strings.Replace(path, ".", "/", -1) + "/detail"
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            uri,
		responseType:   response.ResourceDetailResponse,
		request:        types.NewRequest(nil),
	}
}

func (w *WeCrossRPCModel) Call(path, method string, args ...string) *RemoteCall {
	transactionRequest := request.NewTransactionRequest(method, args)
	currentTxCtx := w.weCrossService.GetTransactionContex()
	txID := currentTxCtx.TxID
	if txID != "" && currentTxCtx.IsPathInTransaction(path) {
		transactionRequest.AddOption(common.XA_TRANSACTION_ID_KEY, txID)
		w.logger.Infof("call: TransactionID exist, turn to callTransaction, TransactionID is %s", txID)
	}
	uri := "/resource/" + strings.Replace(path, ".", "/", -1) + "/call"
	wrappedRequest := types.NewRequest(transactionRequest)
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            uri,
		responseType:   response.TransactionResponse,
		request:        wrappedRequest,
	}
}

func (w *WeCrossRPCModel) SendTransaction(path, method string, args ...string) *RemoteCall {
	transactionRequest := request.NewTransactionRequest(method, args)
	return w.buildSendTransactionRequest(path, transactionRequest)
}

func (w *WeCrossRPCModel) Invoke(path, method string, args ...string) *RemoteCall {
	transactionRequest := request.NewTransactionRequest(method, args)
	currentTxCtx := w.weCrossService.GetTransactionContex()
	xaTransactionID := currentTxCtx.TxID
	if xaTransactionID != "" && currentTxCtx.IsPathInTransaction(path) {
		transactionRequest.AddOption(common.XA_TRANSACTION_ID_KEY, xaTransactionID)
		xaTransactionSeq := currentTxCtx.CurrentXATransactionSeq()
		transactionRequest.AddOption(common.XA_TRANSACTION_SEQ_KEY, xaTransactionSeq)
		w.logger.Infof("invoke: TransactionID exist, turn to execTransaction, TransactionID is %s, Seq is %d", xaTransactionID, xaTransactionSeq)
	}
	return w.buildSendTransactionRequest(path, transactionRequest)
}

func (w *WeCrossRPCModel) CallXA(transactionID, path, method string, args ...string) *RemoteCall {
	transactionRequest := request.NewTransactionRequest(method, args)
	transactionRequest.AddOption(common.XA_TRANSACTION_ID_KEY, transactionID)
	uri := "resource/" + strings.Replace(path, ".", "/", -1) + "/call"
	wrappedRequest := types.NewRequest(transactionRequest)
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            uri,
		responseType:   response.TransactionResponse,
		request:        wrappedRequest,
	}
}

func (w *WeCrossRPCModel) SendXATransaction(transactionID, path, method string, args ...string) *RemoteCall {
	currentTxCtx := w.weCrossService.GetTransactionContex()
	xaTransactionSeq := currentTxCtx.CurrentXATransactionSeq()
	transactionRequest := request.NewTransactionRequest(method, args)
	transactionRequest.AddOption(common.XA_TRANSACTION_ID_KEY, transactionID)
	transactionRequest.AddOption(common.XA_TRANSACTION_SEQ_KEY, xaTransactionSeq)
	return w.buildSendTransactionRequest(path, transactionRequest)
}

func (w *WeCrossRPCModel) StartXATransaction(transactionID string, paths []string) *RemoteCall {
	xaTransactionRequest := request.NewXATransactionRequest(transactionID, paths)
	wrappedRequest := types.NewRequest(xaTransactionRequest)
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            "/xa/startXATransaction",
		responseType:   response.XAResponse,
		request:        wrappedRequest,
	}
}

func (w *WeCrossRPCModel) CommitXATransaction(transactionID string, paths []string) *RemoteCall {
	xaTransactionRequest := request.NewXATransactionRequest(transactionID, paths)
	wrappedRequest := types.NewRequest(xaTransactionRequest)
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            "/xa/commitXATransaction",
		responseType:   response.XAResponse,
		request:        wrappedRequest,
	}
}

func (w *WeCrossRPCModel) RollbackXATransaction(transactionID string, paths []string) *RemoteCall {
	xaTransactionRequest := request.NewXATransactionRequest(transactionID, paths)
	wrappedRequest := types.NewRequest(xaTransactionRequest)
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            "/xa/rollbackXATransaction",
		responseType:   response.XAResponse,
		request:        wrappedRequest,
	}
}

func (w *WeCrossRPCModel) GetXATransaction(transactionID string, paths []string) *RemoteCall {
	xaTransactionRequest := request.NewXATransactionRequest(transactionID, paths)
	wrappedRequest := types.NewRequest(xaTransactionRequest)
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            "/xa/getXATransaction",
		responseType:   response.XAResponse,
		request:        wrappedRequest,
	}
}

func (w *WeCrossRPCModel) CustomCommand(command string, path string, args ...any) *RemoteCall {
	commandRequest := request.NewCommandRequest(command, args)
	wrappedRequest := types.NewRequest(commandRequest)
	uri := "/resource/" + strings.Replace(path, ".", "/", -1) + "/customCommand"
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            uri,
		responseType:   response.CommandResponse,
		request:        wrappedRequest,
	}
}

func (w *WeCrossRPCModel) ListXATransactions(size int) *RemoteCall {
	listXATransactionsRequest := request.NewListXATransactionsRequest(size)
	wrappedRequest := types.NewRequest(listXATransactionsRequest)
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            "/xa/listXATransactions",
		responseType:   response.XATransactionListResponse,
		request:        wrappedRequest,
	}
}

func (w *WeCrossRPCModel) Register(name, password string) (*RemoteCall, *common.WeCrossSDKError) {
	nameOK, _ := regexp.Match(common.USERNAME_PATTERN, []byte(name))
	passwordOK, _ := regexp.Match(common.PASSWORD_PATTERN, []byte(password))
	if !nameOK || !passwordOK {
		return nil, common.NewWeCrossSDKFromString(common.ILLEGAL_SYMBOL, "Invalid username/password, please check your username/password matches the pattern.")
	}
	registerParams, err := w.buildRegisterParams(name, password)
	if err != nil {
		return nil, err
	}

	stringRequest := request.NewStringRequest(registerParams)
	wrappedRequest := types.NewRequest(stringRequest)

	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            "/auth/register",
		responseType:   response.UAResponse,
		request:        wrappedRequest,
	}, nil

}

func (w *WeCrossRPCModel) Login(name, password string) (*RemoteCall, *common.WeCrossSDKError) {
	uaRequest := request.NewUARequest(name, password)
	loginParams, err := w.buildLoginParams(name, password)
	if err != nil {
		return nil, err
	}
	stringRequest := request.NewStringRequest(loginParams)
	wrappedRequest := types.NewRequest(stringRequest)
	wrappedRequest.Ext = uaRequest
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            "/auth/login",
		responseType:   response.UAResponse,
		request:        wrappedRequest,
	}, nil
}

func (w *WeCrossRPCModel) Logout() *RemoteCall {
	uaRequest := request.NewUARequest("", "")
	wrappedRequest := types.NewRequest(uaRequest)
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            "/auth/logout",
		responseType:   response.UAResponse,
		request:        wrappedRequest,
	}
}

func (w *WeCrossRPCModel) AddChainAccount(chainType string, chainAccount account.ChainAccount) *RemoteCall {
	wrappedRequest := types.NewRequest(chainAccount)
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            "/auth/addChainAccount",
		responseType:   response.UAResponse,
		request:        wrappedRequest,
	}
}

func (w *WeCrossRPCModel) SetDefaultAccount(chainType string, chainAccount account.ChainAccount, keyID int) *RemoteCall {
	var wrappedRequest *types.Request
	if chainAccount != nil {
		wrappedRequest = types.NewRequest(chainAccount)
	} else {
		builtChainAccount := &account.CommonAccount{
			KeyID:       keyID,
			AccountType: account.ChainAccountType(chainType),
			IsDefault:   true,
		}
		wrappedRequest = types.NewRequest(builtChainAccount)
	}

	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            "/auth/setDefaultAccount",
		responseType:   response.UAResponse,
		request:        wrappedRequest,
	}
}

func (w *WeCrossRPCModel) GetCurrentTransactionID() string {
	currentTxCtx := w.weCrossService.GetTransactionContex()
	if currentTxCtx.TxID == "" {
		w.logger.Warnf("getCurrentTransactionID: Current TransactionID is null.")
		return ""
	}
	return currentTxCtx.TxID
}

func (w *WeCrossRPCModel) buildSendTransactionRequest(path string, transactionRequest *request.TransactionRequest) *RemoteCall {
	uri := "/resource/" + strings.Replace(path, ".", "/", -1) + "/sendTransaction"
	wrappedRequest := types.NewRequest(transactionRequest)
	return &RemoteCall{
		weCrossService: w.weCrossService,
		httpMethod:     "POST",
		uri:            uri,
		responseType:   response.TransactionResponse,
		request:        wrappedRequest,
	}
}

func (w *WeCrossRPCModel) buildLoginParams(name, password string) (string, *common.WeCrossSDKError) {
	pubResponse, err := w.QueryPub().Send()
	if err != nil {
		return "", err
	}
	authCodeResponse, err := w.QueryAuthCode().Send()
	if err != nil {
		return "", err
	}

	pub := pubResponse.Data.(*response.Pub).Pub
	authCode := authCodeResponse.Data.(*response.AuthCodeReceipt).AuthCode
	confusedPassword := cryptos.Sha256Hex(eles.LoginSalt + password)

	w.logger.Debugf("login username: %s, pub: %s, randomToken: %s", name, pub, authCode.RandomToken)

	loginRequest := &request.LoginRequest{
		UserName:    name,
		PassWord:    confusedPassword,
		RandomToken: authCode.RandomToken,
	}

	pubKey, errr := cryptos.CreatePublicKey(pub)
	if errr != nil {
		return "", common.NewWeCrossSDKFromError(common.RPC_ERROR, errr)
	}

	params, errr := cryptos.EncryptBase64(loginRequest.ToJson(), pubKey)
	if errr != nil {
		return "", common.NewWeCrossSDKFromError(common.INTERNAL_ERROR, errr)
	}
	return params, nil

}

func (w *WeCrossRPCModel) buildRegisterParams(name, password string) (string, *common.WeCrossSDKError) {
	pubResponse, err := w.QueryPub().Send()
	if err != nil {
		return "", err
	}
	authCodeResponse, err := w.QueryAuthCode().Send()
	if err != nil {
		return "", err
	}
	pub := pubResponse.Data.(*response.Pub).Pub
	authCode := authCodeResponse.Data.(*response.AuthCodeReceipt).AuthCode
	confusedPassword := cryptos.Sha256Hex(eles.LoginSalt + password)

	w.logger.Debugf("register username: %s, pub: %s, randomToken: %s", name, pub, authCode.RandomToken)

	registerRequest := new(request.RegisterRequest)
	registerRequest.PassWord = confusedPassword
	registerRequest.UserName = name
	registerRequest.RandomToken = authCode.RandomToken

	pubKey, errr := cryptos.CreatePublicKey(pub)
	if errr != nil {
		return "", common.NewWeCrossSDKFromError(common.RPC_ERROR, errr)
	}
	params, errr := cryptos.EncryptBase64(registerRequest.ToJson(), pubKey)
	if errr != nil {
		return "", common.NewWeCrossSDKFromError(common.INTERNAL_ERROR, errr)
	}
	return params, nil
}

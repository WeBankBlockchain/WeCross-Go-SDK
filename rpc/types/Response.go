package types

import (
	"WeCross-Go-SDK/common"
	"WeCross-Go-SDK/rpc/eles/account"
	"WeCross-Go-SDK/rpc/eles/resource"
	"WeCross-Go-SDK/rpc/types/response"
	"WeCross-Go-SDK/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	Version   string           `json:"version"`
	ErrorCode common.ErrorCode `json:"errorCode"`
	Message   string           `json:"message"`
	Data      Data             `json:"data"`
}

func NewResponse(errorCode common.ErrorCode, message string, data Data) *Response {
	return &Response{
		Version:   common.CURRENT_VERSION,
		ErrorCode: errorCode,
		Message:   message,
		Data:      data,
	}
}

func (rp *Response) ToString() string {
	dataString := ""
	if rp.Data != nil {
		dataString = rp.Data.ToString()
	}
	str := fmt.Sprintf("Response{version='%s', errorCode=%d, message='%s', data=%s}", rp.Version, rp.ErrorCode, rp.Message, dataString)
	return str
}

// ParseResponse use the given responseType to parse a Response
func ParseResponse(httpResponse *http.Response, responseType response.ResponseType) *Response {
	tempResponse := new(Response)
	tempResponse.Version = common.CURRENT_VERSION

	if httpResponse == nil {
		tempResponse.ErrorCode = common.RPC_ERROR
		return tempResponse
	}
	defer httpResponse.Body.Close()

	_, ok := response.ValidResponseTypes[responseType]
	if !ok {
		tempResponse.ErrorCode = common.INTERNAL_ERROR
		return tempResponse
	}

	jsonBytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		tempResponse.ErrorCode = common.RPC_ERROR
		return tempResponse
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(jsonBytes, &m)
	if err != nil {
		tempResponse.ErrorCode = common.RPC_ERROR
		return tempResponse
	}

	tempResponse.Version = utils.AnyToString(m["version"])
	tempResponse.ErrorCode = utils.AnyToErrorCode(m["errorCode"])
	tempResponse.Message = utils.AnyToString(m["message"])

	var data Data
	dataBytes, err := json.Marshal(m["data"])

	switch responseType { // TODO: Add your response type here!
	case response.CommonResponse:
		nullResponse := new(response.NullResponse)
		data = nullResponse
	case response.StubResponse:
		stubs := new(response.Stubs)
		stubs.ParseSelfFromJson(dataBytes)
		data = stubs
	case response.PubResponse:
		pub := new(response.Pub)
		pub.ParseSelfFromJson(dataBytes)
		data = pub
	case response.AuthCodeResponse:
		authCodeReceipt := new(response.AuthCodeReceipt)
		authCodeReceipt.ParseSelfFromJson(dataBytes)
		data = authCodeReceipt
	case response.UAResponse:
		uaReceipt := new(response.UAReceipt)
		uaReceipt.ParseSelfFromJson(dataBytes)
		data = uaReceipt
	case response.XAResponse:
		rawXAResponse := new(response.RawXAResponse)
		rawXAResponse.ParseSelfFromJson(dataBytes)
		data = rawXAResponse
	case response.AccountResponse:
		universalAccount := account.ParseUniversalAccountFromJson(dataBytes)
		data = universalAccount
	case response.ResourceResponse:
		resources := new(resource.Resources)
		resources.ParseSelfFromJson(dataBytes)
		data = resources
	case response.TransactionResponse:
		txReceipt := new(response.TXReceipt)
		txReceipt.ParseSelfFromJson(dataBytes)
		data = txReceipt
	case response.CommandResponse:
		stringResponse := new(response.StringResponse)
		stringResponse.ParseSelfFromJson(dataBytes)
		data = stringResponse
	case response.XATransactionListResponse:
		rawXATransactionListResponse := new(response.RawXATransactionListResponse)
		rawXATransactionListResponse.ParseSelfFromJson(dataBytes)
		data = rawXATransactionListResponse

	default:

	}
	tempResponse.Data = data

	return tempResponse
}

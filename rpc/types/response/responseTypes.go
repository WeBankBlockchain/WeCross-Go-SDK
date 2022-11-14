package response

import (
	"WeCross-Go-SDK/rpc/eles/account"
	"WeCross-Go-SDK/rpc/eles/resources"
	"reflect"
)

type ResponseType string

const (
	CommonResponse            ResponseType = "CommonResponse"
	StubResponse              ResponseType = "StubResponse"
	PubResponse               ResponseType = "PubResponse"
	AuthCodeResponse          ResponseType = "AuthCodeResponse"
	AccountResponse           ResponseType = "AccountResponse"
	ResourceResponse          ResponseType = "ResourceResponse"
	ResourceDetailResponse    ResponseType = "ResourceDetailResponse"
	TransactionResponse       ResponseType = "TransactionResponse"
	CommandResponse           ResponseType = "CommandResponse"
	XATransactionListResponse ResponseType = "XATransactionListResponse"
	UAResponse                ResponseType = "UAResponse"
	XAResponse                ResponseType = "XAResponse"
)

var (
	ValidResponseTypes = map[ResponseType]reflect.Type{
		CommonResponse:            reflect.TypeOf(new(NullResponse)),
		StubResponse:              reflect.TypeOf(new(Stubs)),
		PubResponse:               reflect.TypeOf(new(Pub)),
		AuthCodeResponse:          reflect.TypeOf(new(AuthCodeReceipt)),
		AccountResponse:           reflect.TypeOf(new(account.UniversalAccount)),
		ResourceResponse:          reflect.TypeOf(new(resources.Resources)),
		ResourceDetailResponse:    reflect.TypeOf(new(resources.ResourceDetail)),
		TransactionResponse:       reflect.TypeOf(new(TXReceipt)),
		CommandResponse:           reflect.TypeOf(new(StringResponse)),
		XATransactionListResponse: reflect.TypeOf(new(RawXATransactionListResponse)),
		UAResponse:                reflect.TypeOf(new(UAReceipt)),
		XAResponse:                reflect.TypeOf(new(RawXAResponse)),
	}
)

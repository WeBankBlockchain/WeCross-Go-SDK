package response

import "reflect"

type ResponseType string

const (
	CommonResponse ResponseType = "CommonResponse"
	StubResponse   ResponseType = "StubResponse"
	UAResponse     ResponseType = "UAResponse"
	XAResponse     ResponseType = "XAResponse"
)

var (
	ValidResponseTypes = map[ResponseType]reflect.Type{
		CommonResponse: reflect.TypeOf(new(NullResponse)),
		StubResponse:   reflect.TypeOf(new(Stubs)),
		UAResponse:     reflect.TypeOf(new(UAReceipt)),
		XAResponse:     reflect.TypeOf(new(RawXAResponse)),
	}
)

package response

import "reflect"

type ResponseType string

var (
	ValidResponseTypes = map[ResponseType]reflect.Type{
		"UAResponse": reflect.TypeOf(new(UAReceipt)),
	}
)

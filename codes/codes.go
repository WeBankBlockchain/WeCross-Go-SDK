package codes

import (
	"fmt"
	"strconv"
)

type Code uint32

const (
	Success Code = 0
	// Canceled indicates the operation was canceled (typically by the caller).
	Canceled Code = 1

	InternalError Code = 100

	// config
	FieldMissing  Code = 101
	ResourceError Code = 102
	IllegalSymbol Code = 103

	// rpc
	RemoteCallError    Code = 201
	RpcError           Code = 202
	CallContractError  Code = 203
	LackAuthentication Code = 204

	// performance
	ResourceInactive Code = 301
	InvalidContract  Code = 302

	_maxCode = 9999
)

var strToCode = map[string]Code{
	`"SUCCESS"`:             Success,
	`"INTERNAL_ERROR"`:      InternalError,
	`"FIELD_MISSING"`:       FieldMissing,
	`"RESOURCE_ERROR"`:      ResourceError,
	`"ILLEGAL_SYMBOL"`:      IllegalSymbol,
	`"REMOTE_CALL_ERROR"`:   RemoteCallError,
	`"RPC_ERROR"`:           RpcError,
	`"CALL_CONTRACT_ERROR"`: CallContractError,
	`"LACK_AUTHENTICATION"`: LackAuthentication,
	`"RESOURCE_INACTIVE"`:   ResourceInactive,
	`"INVALID_CONTRACT"`:    InvalidContract,
}

// UnmarshalJSON unmarshals b into the Code.
func (c *Code) UnmarshalJSON(b []byte) error {
	// From json.Unmarshaler: By convention, to approximate the behavior of
	// Unmarshal itself, Unmarshalers implement UnmarshalJSON([]byte("null")) as
	// a no-op.
	if string(b) == "null" {
		return nil
	}
	if c == nil {
		return fmt.Errorf("nil receiver passed to UnmarshalJSON")
	}

	if ci, err := strconv.ParseUint(string(b), 10, 32); err == nil {
		if ci >= _maxCode {
			return fmt.Errorf("invalid code: %q", ci)
		}

		*c = Code(ci)
		return nil
	}

	if jc, ok := strToCode[string(b)]; ok {
		*c = jc
		return nil
	}
	return fmt.Errorf("invalid code: %q", string(b))
}

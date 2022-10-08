package common

type ErrorCode int

const (
	INTERNAL_ERROR ErrorCode = 100

	// config
	FIELD_MISSING  ErrorCode = 101
	RESOURCE_ERROR ErrorCode = 102
	ILLEGAL_SYMBOL ErrorCode = 103

	// rpc
	REMOTECALL_ERROR    ErrorCode = 201
	RPC_ERROR           ErrorCode = 202
	CALL_CONTRACT_ERROR ErrorCode = 203
	LACK_AUTHENTICATION ErrorCode = 204

	// performance
	RESOURCE_INACTIVE ErrorCode = 301
	INVALID_CONTRACT  ErrorCode = 302
)

type StatusCode int

const (
	SUCCESS StatusCode = 0
)

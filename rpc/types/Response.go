package types

import (
	"WeCross-Go-SDK/common"
	"fmt"
)

type Response struct {
	version   string
	errorCode common.ErrorCode
	message   string
	data      Data
}

func NewResponse(errorCode common.ErrorCode, message string, data Data) *Response {
	return &Response{
		version:   common.CURRENT_VERSION,
		errorCode: errorCode,
		message:   message,
		data:      data,
	}
}

func (rp *Response) ToString() string {
	str := fmt.Sprintf("Response{version='%s', errorCode=%d, message='%s', data=%s}", rp.version, rp.errorCode, rp.message, rp.data.ToString())
	return str
}

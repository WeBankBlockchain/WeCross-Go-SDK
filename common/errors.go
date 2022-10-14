package common

import "fmt"

type WeCrossSDKError struct {
	Code    ErrorCode
	Content string
}

func NewWeCrossSDKFromString(errorCode ErrorCode, content string) *WeCrossSDKError {
	return &WeCrossSDKError{
		Code:    errorCode,
		Content: content,
	}
}

func NewWeCrossSDKFromError(errorCode ErrorCode, err error) *WeCrossSDKError {
	return &WeCrossSDKError{
		Code:    errorCode,
		Content: err.Error(),
	}
}

func (w *WeCrossSDKError) Error() string {
	return fmt.Sprintf("WeCrossSDK error with code: %d, %s", w.Code, w.Content)
}

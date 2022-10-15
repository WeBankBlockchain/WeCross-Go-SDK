package utils

import "WeCross-Go-SDK/common"

func AnyToString(obj any) string {
	if obj == nil {
		return ""
	}
	str, _ := obj.(string)
	return str
}

func AnyToInt(obj any) int {
	if obj == nil {
		return 0
	}
	temp, _ := obj.(float64)
	return int(temp)
}

func AnyToErrorCode(obj any) common.ErrorCode {
	return common.ErrorCode(AnyToInt(obj))
}

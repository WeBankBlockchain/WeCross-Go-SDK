package response

import (
	"WeCross-Go-SDK/common"
	"encoding/json"
	"fmt"
)

type AuthCodeReceipt struct {
	ErrorCode common.ErrorCode `json:"errorCode"`
	Message   string           `json:"message"`
	AuthCode  *AuthCodeInfo    `json:"authCode"`
}

func (a *AuthCodeReceipt) ToString() string {
	str := fmt.Sprintf("AuthCodeReceipt{errorCode=%d, message='%s', authCode=%s}", a.ErrorCode, a.Message, a.AuthCode.ToString())
	return str
}

func (a *AuthCodeReceipt) ParseSelfFromJson(valueBytes []byte) {
	err := json.Unmarshal(valueBytes, a)
	if err != nil {
		a.ErrorCode = common.RPC_ERROR

	}
}

type AuthCodeInfo struct {
	RandomToken string `json:"randomToken"`
	ImageBase64 string `json:"imageBase64"`
}

func (a *AuthCodeInfo) ToString() string {
	str := fmt.Sprintf("AuthCodeInfo{randomToken='%s', imageBase64='%s'}", a.RandomToken, a.ImageBase64)
	return str
}

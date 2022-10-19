package request

import (
	"WeCross-Go-SDK/logger"
	"WeCross-Go-SDK/rpc/eles/account"
	"fmt"
)

var UARequestLogger = logger.NewLogger("UARequest")

type UARequest struct {
	UserName     string               `json:"username"`
	PassWord     string               `json:"password"`
	ClientType   string               `json:"clientType"`
	AuthCode     string               `json:"authCode"`
	ChainAccount account.ChainAccount `json:"chainAccount"`
}

func NewUARequest(username, password string) *UARequest {
	return &UARequest{
		UserName:     username,
		PassWord:     password,
		ClientType:   "sdk",
		AuthCode:     "",
		ChainAccount: nil,
	}
}

func (ur *UARequest) ToString() string {
	str := fmt.Sprintf("UARequest{username='%s', password='%s', clientType='%s', authCode='%s'", ur.UserName, ur.PassWord, ur.ClientType, ur.AuthCode)
	if ur.ChainAccount != nil {
		str += fmt.Sprintf(", chainAccount='%s'", ur.ChainAccount.ToString())
	}
	str += "}"
	return str
}

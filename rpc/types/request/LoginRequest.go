package request

import (
	"encoding/json"
	"fmt"
)

type LoginRequest struct {
	UserName    string `json:"username"`
	PassWord    string `json:"password"`
	RandomToken string `json:"randomToken"`
}

func (l *LoginRequest) ToString() string {
	str := fmt.Sprintf("LoginRequest{username='%s', password='%s', randomToken='%s'}", l.UserName, l.PassWord, l.RandomToken)
	return str
}

func (l *LoginRequest) ToJson() []byte {
	jsonBytes, _ := json.Marshal(l)
	return jsonBytes
}

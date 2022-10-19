package request

import (
	"encoding/json"
	"fmt"
)

type RegisterRequest struct {
	UserName    string `json:"username"`
	PassWord    string `json:"password"`
	RandomToken string `json:"randomToken"`
}

func (r *RegisterRequest) ToString() string {
	str := fmt.Sprintf("RegisterRequest{username='%s', password='%s', randomToken='%s'}", r.UserName, r.PassWord, r.RandomToken)
	return str
}

func (r *RegisterRequest) ToJson() []byte {
	jsonBytes, _ := json.Marshal(r)
	return jsonBytes
}

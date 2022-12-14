package types

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/common"
	"encoding/json"
	"fmt"
)

type Request struct {
	Version  string    `json:"version"`
	Data     Data      `json:"data"`
	Ext      Data      `json:"-"`
	CallBack *CallBack `json:"-"`
}

func NewRequest(data Data) *Request {
	return &Request{
		Version: common.CURRENT_VERSION,
		Data:    data,
	}
}

func (req *Request) ToString() string {
	str := fmt.Sprintf("Request{version='%s'", req.Version)
	if req.Data == nil {
		str += ""
	} else {
		str += ", data=" + req.Data.ToString()
	}
	str += "}"
	return str
}

func (req *Request) GetVersion() string {
	return req.Version
}

func (req *Request) ToJson() []byte {
	value, _ := json.Marshal(req)
	return value
}

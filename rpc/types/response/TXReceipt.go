package response

import (
	"WeCross-Go-SDK/common"
	"encoding/json"
	"fmt"
)

type TXReceipt struct {
	ErrorCode   common.ErrorCode `json:"errorCode"`
	Message     string           `json:"message"`
	Hash        string           `json:"hash"`
	ExtraHashes []string         `json:"extraHashes"`
	BlockNumber int              `json:"blockNumber"`
	Result      []string         `json:"result"`
}

func (tr *TXReceipt) ToString() string {
	str := fmt.Sprintf("Receipt{errorCode=%d, errorMessage='%s', hash='%s', extraHashes=%v, blockNumber=%d, result=%v}", tr.ErrorCode, tr.Message, tr.Hash, tr.ExtraHashes, tr.BlockNumber, tr.Result)
	return str
}

func (tr *TXReceipt) ParseSelfFromJson(valueBytes []byte) {
	err := json.Unmarshal(valueBytes, tr)
	if err != nil {
		tr.ErrorCode = common.RPC_ERROR
		tr.Message = err.Error()
	}
}

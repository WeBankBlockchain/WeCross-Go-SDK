package response

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/eles/xa"
)

type RawXATransactionResponse struct {
	XaResponse    *RawXAResponse    `json:"xaResponse"`
	XaTransaction *xa.XATransaction `json:"xaTransaction"`
}

func (xatp *RawXATransactionResponse) ToString() string {
	var rsp, tx string
	if xatp.XaResponse != nil {
		rsp = xatp.XaResponse.ToString()
	}
	if xatp.XaTransaction != nil {
		tx = xatp.XaTransaction.ToString()
	}
	str := fmt.Sprintf("RawXATransactionResponse{xaResponse=%s, xaTransaction=%s}", rsp, tx)
	return str
}

func (xatp *RawXATransactionResponse) ParseSelfFromJson(valueBytes []byte) {
	err := json.Unmarshal(valueBytes, xatp)
	if err != nil {
		xatp.XaResponse = nil
		xatp.XaTransaction = nil
	}
}

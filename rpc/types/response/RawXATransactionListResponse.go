package response

import (
	"WeCross-Go-SDK/rpc/eles/xa"
	"encoding/json"
	"fmt"
	"strings"
)

type RawXATransactionListResponse struct {
	XaList      []*xa.XA       `json:"xaList"`
	NextOffsets map[string]int `json:"nextOffsets"`
	Finished    bool           `json:"finished"`
}

func (rx *RawXATransactionListResponse) ToString() string {
	str := "RawXATransactionListResponse{xaList=["
	for i := 0; i < len(rx.XaList); i++ {
		str += rx.XaList[i].ToString()
		str += ", "
	}
	str = strings.TrimSuffix(str, ", ")
	str += "], " + fmt.Sprintf("nextOffsets=%v, finished=%t}", rx.NextOffsets, rx.Finished)
	return str
}

func (rx *RawXATransactionListResponse) ParseSelfFromJson(valueBytes []byte) {
	err := json.Unmarshal(valueBytes, rx)
	if err != nil {
		rx.XaList = make([]*xa.XA, 0)
		rx.NextOffsets = make(map[string]int)
		rx.Finished = false
	}
}

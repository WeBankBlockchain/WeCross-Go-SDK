package response

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/common"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/eles"
	"encoding/json"
	"strconv"
	"strings"
)

type RawXAResponse struct {
	Status             common.StatusCode         `json:"status"`
	ChainErrorMessages []*eles.ChainErrorMessage `json:"chainErrorMessages"`
}

func (rxar *RawXAResponse) ToString() string {
	str := "rawResponse{status=" + strconv.Itoa(int(rxar.Status)) + ", chainErrorMessages=["
	for _, errMsg := range rxar.ChainErrorMessages {
		str += errMsg.ToString()
		str += ", "
	}
	str = strings.TrimSuffix(str, ", ")
	str += "]}"
	return str
}

func (rxar *RawXAResponse) ParseSelfFromJson(valueBytes []byte) {
	err := json.Unmarshal(valueBytes, rxar)
	if err != nil {
		rxar.Status = common.FAIL
		rxar.ChainErrorMessages = nil
	}
}

package xa

import (
	"fmt"
	"strings"
)

type XATransaction struct {
	XaTransactionID    string               `json:"xaTransactionID"`
	Username           string               `json:"username"`
	Status             string               `json:"status"`
	StartTimestamp     int64                `json:"startTimestamp"`
	CommitTimestamp    int64                `json:"commitTimestamp"`
	RollbackTimestamp  int64                `json:"rollbackTimestamp"`
	Paths              []string             `json:"paths"`
	XaTransactionSteps []*XATransactionStep `json:"xaTransactionSteps"`
}

func (xat *XATransaction) ToString() string {
	paths := "["
	for i := 0; i < len(xat.Paths); i++ {
		paths += xat.Paths[i]
		paths += ","
	}
	paths = strings.TrimSuffix(paths, ",")
	paths += "]"

	steps := "["
	for i := 0; i < len(xat.XaTransactionSteps); i++ {
		steps += xat.XaTransactionSteps[i].ToString()
		steps += ","
	}

	steps = strings.TrimSuffix(steps, ",")
	steps += "]"

	str := fmt.Sprintf("XATransactionResponse{xaTransactionID='%s', username='%s', status='%s', startTimestamp=%d, commitTimestamp=%d, rollbackTimestamp=%d, paths=%s, xaTransactionSteps=%s}", xat.XaTransactionID, xat.Username, xat.Status, xat.StartTimestamp, xat.CommitTimestamp, xat.RollbackTimestamp, paths, steps)
	return str
}

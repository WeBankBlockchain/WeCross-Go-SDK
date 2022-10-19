package request

import "fmt"

type XATransactionRequest struct {
	XaTransactionID string   `json:"xaTransactionID"`
	Paths           []string `json:"paths"`
}

func NewXATransactionRequest(xaTransactionID string, paths []string) *XATransactionRequest {
	return &XATransactionRequest{
		XaTransactionID: xaTransactionID,
		Paths:           paths,
	}
}

func (xaRq *XATransactionRequest) ToString() string {
	str := fmt.Sprintf("RoutineRequest{transactionID='%s', paths=%v}", xaRq.XaTransactionID, xaRq.Paths)
	return str
}

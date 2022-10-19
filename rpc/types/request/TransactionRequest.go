package request

import "fmt"

type TransactionRequest struct {
	Method string   `json:"method"`
	Args   []string `json:"args"`

	// Optional args
	// transactionID, paths, etc...
	Options map[string]interface{} `json:"options"`
}

func NewTransactionRequest(method string, args []string) *TransactionRequest {
	return &TransactionRequest{
		Method:  method,
		Args:    args,
		Options: make(map[string]interface{}),
	}
}

func (tr *TransactionRequest) AddOption(key string, obj any) {
	tr.Options[key] = obj
}

func (tr *TransactionRequest) ToString() string {
	str := fmt.Sprintf("TransactionRequest{method='%s', args=%v, options=%v}", tr.Method, tr.Args, tr.Options)
	return str
}

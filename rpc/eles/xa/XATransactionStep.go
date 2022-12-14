package xa

import "fmt"

type XATransactionStep struct {
	XaTransactionSeq int64  `json:"xaTransactionSeq"`
	Username         string `json:"username"`
	Path             string `json:"path"`
	Timestamp        int64  `json:"timestamp"`
	Method           string `json:"method"`
	Args             string `json:"args"`
}

func (xats *XATransactionStep) ToString() string {
	str := fmt.Sprintf("XATransactionStep{xaTransactionSeq=%d, username='%s', path='%s', timestamp=%d, method='%s', args='%s'}", xats.XaTransactionSeq, xats.Username, xats.Path, xats.Timestamp, xats.Method, xats.Args)
	return str
}

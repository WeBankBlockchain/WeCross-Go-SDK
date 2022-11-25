package xa

import "fmt"

type XA struct {
	XaTransactionID string   `json:"xaTransactionID"`
	UserName        string   `json:"username"`
	Status          string   `json:"status"`
	TimeStamp       int64    `json:"timestamp"`
	Paths           []string `json:"paths"`
}

func (X *XA) ToString() string {
	str := fmt.Sprintf("XA{xaTransactionID='%s', username='%s', status='%s', timestamp=%d, paths=%v}", X.XaTransactionID, X.UserName, X.Status, X.TimeStamp, X.Paths)
	return str
}

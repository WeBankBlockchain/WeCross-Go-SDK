package request

import "fmt"

type ListXATransactionsRequest struct {
	Size    int            `json:"size"`
	Offsets map[string]int `json:"offsets"`
}

func (l *ListXATransactionsRequest) ToString() string {
	str := fmt.Sprintf("ListXATransactionsRequest{size=%d, offsets=%v}", l.Size, l.Offsets)
	return str
}

func NewListXATransactionsRequest(size int) *ListXATransactionsRequest {
	return &ListXATransactionsRequest{
		Size:    size,
		Offsets: make(map[string]int),
	}
}

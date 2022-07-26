package wecross

import (
	"context"
	"testing"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/credentials/insecure"
)

func (s) TestWithTimeout(t *testing.T) {
	conn, err := Dial("http://43.158.200.203:8250",
		WithBlock(),
		WithTransportCredentials(insecure.NewCredentials()))
	if err == nil {
		conn.Close()
	}
	if err != context.DeadlineExceeded {
		t.Fatalf("Dial(_, _) = %v, %v, want %v", conn, err, context.DeadlineExceeded)
	}
}

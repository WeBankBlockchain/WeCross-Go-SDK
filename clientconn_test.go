package wecross

import (
	"context"
	"testing"
	"time"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/credentials/insecure"
)

func (s) TestWithTimeout(t *testing.T) {
	conn, err := Dial("passthrough:///Non-Existent.Server:8250",
		WithTimeout(time.Millisecond),
		WithBlock(),
		WithTransportCredentials(insecure.NewCredentials()))
	if err == nil {
		conn.Close()
	}
	if err != context.DeadlineExceeded {
		t.Fatalf("Dial(_, _) = %v, %v, want %v", conn, err, context.DeadlineExceeded)
	}
}

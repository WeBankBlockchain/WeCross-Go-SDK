package wecross

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/transport"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/credentials/insecure"
)

type failFastError struct{}

func (failFastError) Error() string   { return "falifast" }
func (failFastError) Temporary() bool { return false }

func TestDialContextFailFast(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
	defer cancel()
	failErr := failFastError{}
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return nil, failErr
	}

	_, err := DialContext(ctx, "Non-Existent.Server:8250", WithBlock(), WithTransportCredentials(insecure.NewCredentials()), WithContextDialer(dialer), FailOnNonTempDialError(true))
	if terr, ok := err.(transport.ConnectionError); !ok || terr.Origin() != failErr {
		t.Fatalf("DialContext() = _, %v, want _, %v", err, failErr)
	}
}

func (s) TestWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	conn, err := DialContext(ctx, "passthrough:///Non-Existent.Server:8250",
		WithBlock(),
		WithTransportCredentials(insecure.NewCredentials()))
	if err == nil {
		conn.Close()
	}
	if err != context.DeadlineExceeded {
		t.Fatalf("Dial(_, _) = %v, %v, want %v", conn, err, context.DeadlineExceeded)
	}
}

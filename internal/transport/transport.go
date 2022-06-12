package transport

import "github.com/WeBankBlockchain/WeCross-Go-SDK/credentials"

// ConnectOptions covers all relevant options for coummunicating with the server.
type ConnectOptions struct {
	// FailOnNonTempDialError specifies if WeCross fails on non-temporary dial errors.
	FailOnNonTempDialError bool
	// PerRPCCredentials stores the PerRPCCredentials required to issue RPCs.
	// Note that PerPRCCredentials has not been implemented.
	PerRPCCredentials []credentials.PerRPCCredentials
	// TransportCredentials stores the Authenticator required to setup a client
	// connection.
	TransportCredentials credentials.TransportCredentials
	CredsBundle          credentials.Bundle
	// WriteBufferSize sets the size of write buffer which in turn determines how much data can be batched before it's written on the wire.
	WriteBufferSize int
	// ReadBufferSize sets the size of read buffer, which in turn determines how much data can be read at most for one read syscall.
	ReadBufferSize int
}

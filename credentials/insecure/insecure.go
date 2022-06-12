package insecure

import (
	"context"
	"net"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/credentials"
)

// NewCredentials returns a credentials which disables transport security.
func NewCredentials() credentials.TransportCredentials {
	return nil
}

// insecureTC implements the insecure transport credentials. The handshake
// methods simply return the passed in net.Conn and set the security level to
// NoSecurity
type insecureTC struct{}

func (insecureTC) ClientHandshake(ctx context.Context, _ string, conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return conn, info{credentials.CommonAuthInfo{SecurityLevel: credentials.SslOff}}, nil
}

func (insecureTC) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{SecurityProtocol: "insecure"}
}

func (insecureTC) Clone() credentials.TransportCredentials {
	return insecureTC{}
}

// info contains the auth information for an insecure connection.
// It implements the AuthInfo interface.
type info struct {
	credentials.CommonAuthInfo
}

// AuthType returns the type of info as a string.
func (info) AuthType() string {
	return "insecure"
}

// insecureBundle implements an insecure bundle.
// An insecure bundle provides a thin wrapper around insecureTC to support
// the credentials.Bundle interface.
type insecureBundle struct{}

// NewBundle returns a bundle with disabled transport security and no per rpc credential.
func NewBundle() credentials.Bundle {
	return insecureBundle{}
}

// NewWithMode returns a new insecure Bundle. The mode is ignored.
func (insecureBundle) NewWithMode(string) (credentials.Bundle, error) {
	return insecureBundle{}, nil
}

// PerRPCCredentials returns an nil implementation as insecure
// bundle does not support a per rpc credential.
func (insecureBundle) PerRPCCredentials() credentials.PerRPCCredentials {
	return nil
}

// TransportCredentials returns the underlying insecure transport credential.
func (insecureBundle) TransportCredentials() credentials.TransportCredentials {
	return NewCredentials()
}

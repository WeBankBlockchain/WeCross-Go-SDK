package credentials

import (
	"context"
	"fmt"
	"net"
)

// PerRPCCredentials defines the common interface for the credentials which need to
// attach security information to every RPC.
//
// Note that PerPRCCredentials has not been implemented.
type PerRPCCredentials interface {
	// GetRequestMetadata gets the current request metadata, refreshing
	// tokens if required. This should be called by the transport layer on
	// each request, and the data should be populated in headers or other
	// context. If a status code is returned, it will be used as the status
	// for the RPC. uri is the URI of the entry point for the request.
	// When supported by the underlying implementation, ctx can be used for
	// timeout and cancellation. Additionally, RequestInfo data will be
	// available via ctx to this call.
	// it as an arbitrary string.
	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
	// RequireTransportSecurity indicates whether the credentials requires
	// transport security.
	RequireTransportSecurity() bool
}

// AuthInfo defines the common interface for the auth information the users are interested in.
// A struct that implements AuthInfo should embed CommonAuthInfo by including additional
// information about the credentials in it.
type AuthInfo interface {
	AuthType() string
}

// SecurityLevel defines the protection level on an established connection.
type SecurityLevel int

const (
	SslOnClientAuth SecurityLevel = iota
	SslOff
	SslOn
)

// String returns SecurityLevel in a string format.
func (s SecurityLevel) String() string {
	switch s {
	case SslOnClientAuth:
		return "SslOnClientAuth"
	case SslOff:
		return "SslOff"
	case SslOn:
		return "SslOn"
	}
	return fmt.Sprintf("invalid SecurityLevel: %v", int(s))
}

// CommonAuthInfo contains authenticated information common to AuthInfo implementations.
// It should be embedded in a struct implementing AuthInfo to provide additional information
// about the credentials.
type CommonAuthInfo struct {
	SecurityLevel SecurityLevel
}

// ProtocolInfo provides information regarding the gRPC wire protocol version,
// security protocol, security protocol version in use, server name, etc.
type ProtocolInfo struct {
	// ProtocolVersion is the wecross wire protocol version.
	ProtocolVersion string
	// SecurityProtocol is the security protocol in use.
	SecurityProtocol string
}

// TransportCredentials defines the common interface for all the live weCross
// protocols and supported transport security protocols (e.g., SSL).
type TransportCredentials interface {
	// ClientHandshake does the authentication handshake specified by the
	// corresponding authentication protocol on rawConn for clients. It returns
	// the authenticated connection and the corresponding auth information
	// about the connection. The auth information should embed CommonAuthInfo
	// to return additional information about the credentials. Implementations
	// must use the provided context to implement timely cancellation.
	// Additionally, ClientHandshakeInfo data will be available via the context
	// passed to this call.
	//
	// The second argument to this method is the `host:port` header value used
	// while creating new streams on this connection after authentication
	// succeeds.
	//
	// If the returned net.Conn is closed, it MUST close the net.Conn provided.
	ClientHandshake(context.Context, string, net.Conn) (net.Conn, AuthInfo, error)
	// Clone makes a copy of this TransportCredentials.
	Clone() TransportCredentials
	// Info provides the ProtocolInfo of this TransportCredentials.
	Info() ProtocolInfo
}

// Bundle is a combination of TransportCredentials and PerRPCCredentials.
//
// It also contains a mode switching method, so it can be used as a combination
// of different credential policies.
//
// Bundle cannot be used together with individual TransportCredentials.
// PerPRCCredentials from Bundle will be appended to other PerRPCCredentials.
//
// Note that PerPRCCredentials has not been implemented.
type Bundle interface {
	// TransportCredentials returns the transport credentials from the Bundle.
	//
	// Implementations must return non-nil transport credentials. If transport
	// security is not needed by the Bundle, implementations may choose to
	// return insecure.NewCredentials().
	TransportCredentials() TransportCredentials

	// PerRPCCredentials returns the per-RPC credentials from the Bundle.
	//
	// May be nil if per-RPC credentials are not needed.
	// Note that PerPRCCredentials has not been implemented.
	PerRPCCredentials() PerRPCCredentials

	// NewWithMode should make a copy of Bundle, and switch mode. Modifying the
	// existing Bundle may cause races.
	//
	// NewWithMode returns nil if the requested mode is not supported.
	NewWithMode(mode string) (Bundle, error)
}

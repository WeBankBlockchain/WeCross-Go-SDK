package credentials

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"

	credinternal "github.com/WeBankBlockchain/WeCross-Go-SDK/internal/credentials"
)

type TLSInfo struct {
	State tls.ConnectionState
	CommonAuthInfo
	SPIFFEID *url.URL
}

// AuthType returns the type of TLSInfo as a string.
func (t TLSInfo) AuthType() string {
	return "tls"
}

// tlsCreds is the credentials required for authenticating a connection using TLS.
type tlsCreds struct {
	// TLS configuration
	config *tls.Config
}

func (c tlsCreds) Info() ProtocolInfo {
	return ProtocolInfo{
		SecurityProtocol: "tls",
	}
}

func (c *tlsCreds) ClientHandshake(ctx context.Context, authority string, rawConn net.Conn) (_ net.Conn, _ AuthInfo, err error) {
	cfg := credinternal.CloneTLSConfig(c.config)
	if cfg.ServerName == "" {
		serverName, _, err := net.SplitHostPort(authority)
		if err != nil {
			// If the authority had no host port or if the authority cannot be parsed, use it as-is.
			serverName = authority
		}
		cfg.ServerName = serverName
	}
	conn := tls.Client(rawConn, cfg)
	errChannel := make(chan error, 1)
	go func() {
		errChannel <- conn.Handshake()
		close(errChannel)
	}()
	select {
	case err := <-errChannel:
		if err != nil {
			conn.Close()
			return nil, nil, err
		}
	case <-ctx.Done():
		conn.Close()
		return nil, nil, ctx.Err()
	}
	tlsInfo := TLSInfo{
		State: conn.ConnectionState(),
		CommonAuthInfo: CommonAuthInfo{
			SecurityLevel: SslOnClientAuth,
		},
	}
	id := credinternal.SPIFFEIDFromState(conn.ConnectionState())
	if id != nil {
		tlsInfo.SPIFFEID = id
	}
	return credinternal.WrapSyscallConn(rawConn, conn), tlsInfo, nil
}

func (c *tlsCreds) Clone() TransportCredentials {
	return NewTLS(c.config)
}

// NewTLS uses c to construct a TransportCredentials based on TLS.
func NewTLS(c *tls.Config) TransportCredentials {
	tc := &tlsCreds{credinternal.CloneTLSConfig(c)}
	tc.config.NextProtos = credinternal.AppendH2ToNextProtos(tc.config.NextProtos)
	return tc
}

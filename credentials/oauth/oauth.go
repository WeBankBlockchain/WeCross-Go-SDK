package oauth

import (
	"context"
)

// Note that PerPRCCredentials has not been implemented.
type jwtAccess struct {
	jsonKey []byte
}

func (j jwtAccess) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return nil, nil
}

func (j jwtAccess) RequireTransportSecurity() bool {
	return true
}

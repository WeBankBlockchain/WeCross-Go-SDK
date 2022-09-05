package transport

import (
	"math"
	"time"
)

const (
	infinity                      = time.Duration(math.MaxInt64)
	defaultClientKeepaliveTime    = infinity
	defaultClientKeppaliveTimeout = 20 * time.Second
)

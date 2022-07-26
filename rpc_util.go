package wecross

// callInfo contains all related configuration and information about an RPC.
type callInfo struct {
	failFast              bool
	maxReceiveMessageSize *int
	maxSendMessageSize    *int
	maxRetryRPCBufferSize int
}

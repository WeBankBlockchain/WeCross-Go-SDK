package eles

var AUTH_REQUIRED_COMMANDS = map[string]bool{
	"call":                  true,
	"sendTransaction":       true,
	"invoke":                true,
	"callXA":                true,
	"sendXATransaction":     true,
	"startXATransaction":    true,
	"commitXATransaction":   true,
	"rollbackXATransaction": true,
	"getXATransaction":      true,
	"listXATransactions":    true,
	"logout":                true,
	"listAccounts":          true,
	"customCommand":         true,
	"addChainAccount":       true,
	"setDefaultAccount":     true,
	"listAccount":           true,
	"supportedStubs":        true,
	"listResources":         true,
	"status":                true,
	"detail":                true,
}

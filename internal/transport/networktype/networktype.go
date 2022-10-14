package networktype

import "github.com/WeBankBlockchain/WeCross-Go-SDK/resolver"

// keyType is the key to use for storing State in Attributes.
type keyType string

const key = keyType("wecross.internal.transport.networktypes")

// Set returns a copy of the provided address with attributes containing networkType.
func Set(address resolver.Address, networkType string) resolver.Address {
	address.Attributes = address.Attributes.WithValue(key, networkType)
	return address
}

// Get returns the network type in the resolver.Address and true, or "", false
// if not present.
func Get(address resolver.Address) (string, bool) {
	v := address.Attributes.Value(key)
	if v == nil {
		return "", false
	}
	return v.(string), true
}

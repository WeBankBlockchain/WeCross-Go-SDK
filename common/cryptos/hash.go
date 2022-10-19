package cryptos

import (
	"crypto/sha256"
	"encoding/hex"
)

// Sha256Hex computes out a hash of sha256 and return it as a hex string
func Sha256Hex(src string) string {
	hash := sha256.New()
	hash.Write([]byte(src))
	res := hex.EncodeToString(hash.Sum(nil))
	return res
}

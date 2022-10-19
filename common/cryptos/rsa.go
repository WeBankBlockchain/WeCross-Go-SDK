package cryptos

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

func CreatePublicKey(content string) (*rsa.PublicKey, error) {
	plainCotent, err := base64Decode(content)
	if err != nil {
		return nil, err
	}
	pubKey, err := x509.ParsePKIXPublicKey(plainCotent)
	if err != nil {
		return nil, err
	}
	x509KeySpec, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot convert to rsa public key")
	}
	return x509KeySpec, nil
}

func EncryptBase64(src []byte, pubKey *rsa.PublicKey) (string, error) {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, src)
	if err != nil {
		return "", err
	}
	plainEncrypted := base64Encode(encrypted)
	return plainEncrypted, nil
}

func base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func base64Decode(str string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return data, nil
}

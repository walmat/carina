package adyen2

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"math/big"
)

type adrsa struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func (r *adrsa) SetPublicKey(pubHex string, e int) error {
	k := &rsa.PublicKey{}
	dec, err := hex.DecodeString(pubHex)
	if err != nil {
		return err
	}
	k.N = new(big.Int).SetBytes(dec)
	k.E = e
	r.publicKey = k
	return nil
}

func (r *adrsa) Encrypt(str string, encoding string) (string, error) {
	bytes := []byte(str)
	res, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, bytes)
	if err != nil {
		return "", err
	}
	if encoding == "base64" {
		return base64.StdEncoding.EncodeToString(res), nil
	}
	return hex.EncodeToString(res), nil
}

func (r *adrsa) Encrypt2(bytes []byte, encoding string) (string, error) {
	res, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, bytes)
	if err != nil {
		return "", err
	}
	if encoding == "base64" {
		return base64.StdEncoding.EncodeToString(res), nil
	}
	return hex.EncodeToString(res), nil
}
func NewRsa() *adrsa {
	rs := &adrsa{}
	return rs
}

package adyen

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

func (r *adrsa) Init(pubHex string, e int) error {
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
func (r *adrsa) encryptWithAesKey(bytes []byte) (string, error) {
	encryptedPublicKey, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedPublicKey), nil
}
func NewRsa() *adrsa {
	rs := &adrsa{}
	return rs
}

package apiutil

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"nebula/pkg/api/model"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
	"golang.org/x/crypto/argon2"
)

func ShouldIgnoreCrypto(req *resty.Request) bool {
	_, ok := req.Header[model.IgnoreCryptoHeader]
	return ok
}

func RsaEncrypt(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		data,
		nil)
	if err != nil {
		return nil, err
	}

	return encryptedBytes, nil
}

func ShaHashSum(msg []byte) ([]byte, error) {
	msgHash := sha256.New()
	_, err := msgHash.Write(msg)
	if err != nil {
		return nil, err
	}
	return msgHash.Sum(nil), nil
}

func RsaSign(privateKey *rsa.PrivateKey, msg []byte) ([]byte, error) {
	msgHashSum, err := ShaHashSum(msg)
	if err != nil {
		return nil, err
	}

	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func RsaVerify(publicKey *rsa.PublicKey, msgHashSum, signature []byte) bool {
	err := rsa.VerifyPSS(publicKey, crypto.SHA256, msgHashSum, signature, nil)
	return err == nil
}

var (
	DefaultPasswordHashParams = &PasswordHashParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
)

type PasswordHashParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

/*
	Passwords should be between 8-64 characters.
	A maximum limit is imposed to avoid a string
	length Denial-of-Service attack.

	Password strength isn't taken into consideration
*/
func IsValidPassword(password string) error {
	if len(password) < 8 || len(password) > 64 {
		return errors.New("password must be between 8-64 characters")
	}

	alphanumericRegex := regexp.MustCompile(`([A-Za-z])([0-9])`)
	if !alphanumericRegex.MatchString(password) {
		return errors.New("password must contain a letter and number")
	}

	return nil
}

func HashPassword(password string, p *PasswordHashParams) (string, error) {
	passwordErr := IsValidPassword(password)
	if passwordErr != nil {
		return "", passwordErr
	}

	if p == nil {
		p = DefaultPasswordHashParams
	}

	salt, err := GenerateRandomBytes(p.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func decodeHash(encodedHash string) (*PasswordHashParams, []byte, []byte, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	params := &PasswordHashParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLength = uint32(len(hash))

	return params, salt, hash, nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

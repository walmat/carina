package adyen

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type Adyen struct {
	rsa      *adrsa
	aesKey   []byte
	aesNonce []byte

	prefix string
}

type Data struct {
	Number         string `json:"number"`
	Cvc            string `json:"cvc"`
	ExpiryMonth    string `json:"expiryMonth"`
	ExpiryYear     string `json:"expiryYear"`
	GenerationTime string `json:"generationtime"`
}

type ExpiryYear struct {
	ExpiryYear     string `json:"expiryYear"`
	GenerationTime string `json:"generationtime"`
}

type ExpiryMonth struct {
	ExpiryMonth    string `json:"expiryMonth"`
	GenerationTime string `json:"generationtime"`
}

type CVC struct {
	CVC            string `json:"cvc"`
	GenerationTime string `json:"generationtime"`
}

type CreditCardNumber struct {
	Number         string `json:"number"`
	GenerationTime string `json:"generationtime"`
}

type Info struct {
	EncryptedCardNumber string `json:"encryptedCardNumber"`
	EncryptedExpMonth   string `json:"encryptedExpMonth"`
	EncryptedExpYear    string `json:"encryptedExpYear"`
	EncryptedCvc        string `json:"encryptedCvc"`
}

func NewAdyen(ver string, publicKey string) *Adyen {
	y := &Adyen{
		rsa:    NewRsa(),
		aesKey: make([]byte, 32),
		prefix: ver,
	}

	err := y.rsa.Init(publicKey, 65537)
	if err != nil {
		panic(err)
	}
	return y
}

func (y *Adyen) random(len int) []byte {
	ak := make([]byte, len)
	_, _ = rand.Read(ak)
	return ak
}

type DetailType string

const (
	DetailTypeCardNumber DetailType = "cardNumber"
	DetailTypeExpMonth   DetailType = "expMonth"
	DetailTypeExpYear    DetailType = "expYear"
	DetailTypeCvc        DetailType = "cvc"
)

func (y *Adyen) EncryptCreditcardDetails(cardNumber string, expMonth string, expYear string, cvc string) (*Info, error) {
	encryptedCard, err := y.EncryptDetail(DetailTypeCardNumber, cardNumber)
	if err != nil {
		return nil, err
	}
	encryptedExpMonth, err := y.EncryptDetail(DetailTypeExpMonth, expMonth)
	if err != nil {
		return nil, err
	}
	encryptedExpYear, err := y.EncryptDetail(DetailTypeExpYear, expYear)
	if err != nil {
		return nil, err
	}
	encryptedCvc, err := y.EncryptDetail(DetailTypeCvc, cvc)
	if err != nil {
		return nil, err
	}
	return &Info{
		EncryptedCardNumber: encryptedCard,
		EncryptedExpMonth:   encryptedExpMonth,
		EncryptedExpYear:    encryptedExpYear,
		EncryptedCvc:        encryptedCvc,
	}, nil
}
func (y *Adyen) EncryptDetail(detailType DetailType, data string) (string, error) {
	y.aesKey = y.random(32)
	y.aesNonce = y.random(12)
	gt := time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00")

	var bytes []byte
	switch detailType {
	case DetailTypeCardNumber:
		info := CreditCardNumber{
			Number:         data,
			GenerationTime: gt,
		}
		bytes, _ = json.Marshal(info)
	case DetailTypeExpMonth:
		info := ExpiryMonth{
			ExpiryMonth:    data,
			GenerationTime: gt,
		}
		bytes, _ = json.Marshal(info)
	case DetailTypeExpYear:
		info := ExpiryYear{
			ExpiryYear:     data,
			GenerationTime: gt,
		}
		bytes, _ = json.Marshal(info)
	case DetailTypeCvc:
		info := CVC{
			CVC:            data,
			GenerationTime: gt,
		}
		bytes, _ = json.Marshal(info)
	}

	y.aesKey = y.random(32)
	y.aesNonce = y.random(12)
	block, err := aes.NewCipher(y.aesKey)
	if err != nil {
		return "", err
	}
	cmer, err := NewCCM(block, 8, len(y.aesNonce))
	if err != nil {
		return "", err
	}

	cipherBytes := cmer.Seal(nil, y.aesNonce, bytes, nil)
	cipherBytes = append(y.aesNonce, cipherBytes...)
	cipherText := base64.StdEncoding.EncodeToString(cipherBytes)

	encryptedPublicKey, err := y.rsa.encryptWithAesKey(y.aesKey)
	if err != nil {
		return "", err
	}
	// prefix := "adyenjs_0_1_18$"
	arr := []string{y.prefix, encryptedPublicKey, "$", cipherText}
	return strings.Join(arr, ""), nil
}

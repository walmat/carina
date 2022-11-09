package adyen2

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pion/dtls/v2/pkg/crypto/ccm"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

type Adyen struct {
	rsa     *adrsa
	prefix  string
	version string

	rsaEncryptAesVal string
	//iv
	aesKey     []byte
	_tagSize   int
	_nonceSize int
	_debug     bool
	aesNonce   []byte
}

type Data struct {
	Number              string `json:"number"`
	Cvc                 string `json:"cvc"`
	HolderName          string `json:"holderName"`
	ExpiryMonth         string `json:"expiryMonth"`
	ExpiryYear          string `json:"expiryYear"`
	Generationtime      string `json:"generationtime"`
	PaymentMethodID     string `json:"paymentMethodId"`
	CardType            string `json:"cardType"`
	Referrer            string `json:"referrer"`
	Activate            string `json:"activate"`
	Deactivate          string `json:"deactivate"`
	InitializeCount     string `json:"initializeCount"`
	LuhnCount           string `json:"luhnCount"`
	LuhnOkCount         string `json:"luhnOkCount"`
	LuhnSameLengthCount string `json:"luhnSameLengthCount"`
	SjclStrength        string `json:"sjclStrength"`
	DfValue             string `json:"dfValue"`
}

type CardForm struct {
	Number          string `json:"number" validate:"required"`
	Holder          string `json:"holder" validate:"required"`
	Cvv             string `json:"cvv" validate:"required"`
	ExpiryMonth     string `json:"expiry_month" validate:"required"`
	ExpiryYear      string `json:"expiry_year" validate:"required"`
	CardType        string `json:"card_type" validate:"required"`
	PaymentMethodId string `json:"payment_method_id" validate:"required"`
	Referer         string `json:"referrer"`
	DfVal           string `json:"df_val"` //理论上只指纹数据
}

func NewAdYen(publicKey string) *Adyen {
	yen := &Adyen{}
	yen.rsa = NewRsa()
	yen.prefix = "adyenjs_"
	yen.version = "0_1_21"
	yen.aesKey = make([]byte, 32)
	yen._tagSize = 8
	yen._nonceSize = 12

	//如果密钥错误直接推出
	err := yen.rsa.SetPublicKey(publicKey, 65537)
	if err != nil {
		panic(err)
	}
	return yen
}

func (yen *Adyen) marshal(data interface{}) []byte {
	if reflect.TypeOf(data).String() == "string" {
		return []byte(data.(string))
	}
	bytes, _ := json.Marshal(data)
	return bytes
}

//validate order info
func (yen *Adyen) validate(card CardForm) (Data, error) {
	gt := time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00")
	info := Data{
		Number:              card.Number,
		Cvc:                 card.Cvv,
		HolderName:          card.Holder,
		ExpiryMonth:         card.ExpiryMonth,
		ExpiryYear:          card.ExpiryYear,
		Generationtime:      gt,
		PaymentMethodID:     card.PaymentMethodId,
		CardType:            card.CardType,
		Referrer:            card.Referer,
		Activate:            "4",
		Deactivate:          "4",
		InitializeCount:     "1",
		LuhnCount:           "1",
		LuhnOkCount:         "1",
		LuhnSameLengthCount: "1",
		SjclStrength:        "10",
		DfValue:             "",
	}
	//check card info
	if err := validateAdYenCardInfo(info); err != nil {
		return info, err
	}

	return info, nil
}

//encrypt
func (yen *Adyen) Encrypt(card CardForm) (string, error) {
	// pre validate
	data, err := yen.validate(card)
	bytes := yen.marshal(data)

	//1. init aes random key
	yen.init()

	//2. create ccm instance
	block, err := aes.NewCipher(yen.aesKey)
	if err != nil {
		return "", err
	}
	cmer, err := ccm.NewCCM(block, yen._tagSize, len(yen.aesNonce))
	if err != nil {
		return "", err
	}

	//3. aes encrypt data
	cipherBytes := cmer.Seal(nil, yen.aesNonce, bytes, nil)
	cipherBytes = append(yen.aesNonce, cipherBytes...) //追加
	cipherText := base64.StdEncoding.EncodeToString(cipherBytes)
	if yen._debug {
		fmt.Println("aes:", cipherText)
	}
	//4. rsa encrypt aes.key
	rsaCp, err := yen.rsa.Encrypt2(yen.aesKey, "base64")
	if err != nil {
		return "", err
	}
	if yen._debug {
		fmt.Println("rsa:", rsaCp)
		fmt.Println("rsa:", len(rsaCp))
	}

	//5. append data
	prefix := yen.prefix + yen.version + "$"
	arr := []string{prefix, rsaCp, "$", cipherText}
	return strings.Join(arr, ""), nil
}

func (yen *Adyen) debug() {
	yen._debug = true
	yen.aesKey = []byte{
		211, 230, 56, 196, 255, 13, 107, 44,
		124, 11, 172, 57, 108, 47, 222, 207,
		139, 212, 162, 56, 51, 163, 147, 100,
		195, 176, 241, 192, 75, 86, 32, 68,
	}
	yen.aesNonce = []byte{
		32, 103, 85, 109, 226,
		169, 214, 87, 67, 166,
		9, 92,
	}
}
func (yen *Adyen) init() {
	if yen._debug {
		return
	}
	yen.aesKey = yen.random(32)
	yen.aesNonce = yen.random(12)

}
func (yen *Adyen) random(len int) []byte {
	ak := make([]byte, len)
	_, _ = rand.Read(ak)
	return ak
}

func validateAdYenCardInfo(data Data) error {
	return nil
}

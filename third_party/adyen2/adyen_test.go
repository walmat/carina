package adyen2

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/pion/dtls/v2/pkg/crypto/ccm"
	"testing"
)

func TestAdyen_Encrypt2(t *testing.T) {
	key := []byte{
		175, 152, 214, 174, 41, 51, 112, 162,
		103, 202, 35, 202, 184, 85, 102, 148,
		69, 185, 161, 146, 38, 106, 186, 74,
		49, 210, 186, 58, 7, 43, 219, 17,
	}
	//str := `1`
	//res := common.EcbEncrypt([]byte(str), key)
	//fmt.Println(base64.StdEncoding.EncodeToString(res))
	blk, err := aes.NewCipher(key)
	fmt.Println(blk, err)

	ak := make([]byte, 12)
	_, err = rand.Read(ak)
	ak = []byte{
		78, 6, 182, 161, 114,
		39, 249, 221, 6, 179,
		155, 169}

	plaintext := []byte("1")
	_, err = aes.NewCipher(key)
	fmt.Println(err)
	ccm, err := ccm.NewCCM(blk, 8, 12)
	fmt.Println(err)
	res := ccm.Seal(nil, ak, plaintext, nil)

	res = append(ak, res...)
	fmt.Println(base64.StdEncoding.EncodeToString(res))
}

//原生加密
func TestAdyen_Encrypt(t *testing.T) {
	const rsaKey = "C1B003094C0E0F486D1C88A69756064B3F91F064D908703698357B58FEA57B301C766E73246B9228F4D2ED1B70ABFE3B69DD0592239A90345E450EF95FF6A9B70A7AF2765D336A310F1A293C4DA00ADCC959CD17EA9E93C71FB4DA1655BF99D23CC26FE832A7B61BE9DE3E7E3217D956FFC0565F0B9C4695D73E648F35CFC9853F1EF2AA6A235C6FB2C95BB1B10F0ECC571FE7764ACAFE2002041F0C9B74DE3C375E0FDD6B02B59EC71E7AB2E7242E747A7F8048CDB57120103A7532BF127676BF592804F286DC01C2E79E8CEA44A188F3D4D012774271C755BDC6FD4E5C6D4C59752A2923DA5F9A94E05ABA40CAD65AFDE331481140C4FDE2B314E92FFA3F17"
	yen := NewAdYen(rsaKey)
	yen.debug()
	res, err := yen.Encrypt(CardForm{
		Number:          "xx",
		Holder:          "xx don",
		Cvv:             "x",
		ExpiryMonth:     "x",
		ExpiryYear:      "x",
		CardType:        "VISA",
		PaymentMethodId: "CREDIT_CARD",
		Referer:         "https://www.yeezysupply.com/payment",
	})
	fmt.Println(res, "err:", err)
}

func BenchmarkAdyen_Encrypt(b *testing.B) {
	const rsaKey = "C1B003094C0E0F486D1C88A69756064B3F91F064D908703698357B58FEA57B301C766E73246B9228F4D2ED1B70ABFE3B69DD0592239A90345E450EF95FF6A9B70A7AF2765D336A310F1A293C4DA00ADCC959CD17EA9E93C71FB4DA1655BF99D23CC26FE832A7B61BE9DE3E7E3217D956FFC0565F0B9C4695D73E648F35CFC9853F1EF2AA6A235C6FB2C95BB1B10F0ECC571FE7764ACAFE2002041F0C9B74DE3C375E0FDD6B02B59EC71E7AB2E7242E747A7F8048CDB57120103A7532BF127676BF592804F286DC01C2E79E8CEA44A188F3D4D012774271C755BDC6FD4E5C6D4C59752A2923DA5F9A94E05ABA40CAD65AFDE331481140C4FDE2B314E92FFA3F17"
	yen := NewAdYen(rsaKey)
	for v := 0; v < b.N; v++ {
		_, _ = yen.Encrypt(CardForm{
			Number:          "x",
			Holder:          "iwan don",
			Cvv:             "x",
			ExpiryMonth:     "x",
			ExpiryYear:      "xx",
			CardType:        "VISA",
			PaymentMethodId: "CREDIT_CARD",
			Referer:         "https://www.yeezysupply.com/payment",
		})
		//fmt.Println(res, "err:", err)
	}
}

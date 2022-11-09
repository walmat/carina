package client

import (
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"nebula/pkg/api/apiutil"
	"nebula/pkg/api/model"
)

// ApiUrl injected at build time
var ApiUrl string

var (
	client = resty.New()
)

func InitRSA(pubKey *rsa.PublicKey) {
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		if apiutil.ShouldIgnoreCrypto(req) || req.Method != resty.MethodPost {
			return nil
		}

		if marshalledBody, err := json.Marshal(req.Body); err != nil {
			return err
		} else if encryptedBody, err := apiutil.RsaEncrypt(pubKey, marshalledBody); err != nil {
			return err
		} else if jsonBody, err := json.Marshal(encryptedBody); err != nil {
			return err
		} else {
			req.Body = jsonBody
		}

		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		if apiutil.ShouldIgnoreCrypto(resp.Request) || resp.StatusCode() == fiber.StatusInternalServerError {
			return nil
		} else if len(resp.Body()) == 0 {
			return nil
		}

		encodedSig := resp.Header().Get(model.SignatureHeader)
		if encodedSig == "" {
			return errors.New("no signature in response")
		}

		signature, err := hex.DecodeString(encodedSig)
		if err != nil {
			return err
		}

		msgHashSum, err := apiutil.ShaHashSum(resp.Body())
		if err != nil {
			return err
		}

		if !apiutil.RsaVerify(pubKey, msgHashSum, signature) {
			return errors.New("response does not match signature")
		}

		return nil
	})
}

func Client() *resty.Client {
	return client
}

func R() *resty.Request {
	return client.R()
}

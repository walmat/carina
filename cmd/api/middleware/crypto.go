package middleware

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"nebula/pkg/api/apiutil"
	"nebula/pkg/api/model"
)

func DecryptBody(privateKey *rsa.PrivateKey) iris.Handler {
	return func(ctx iris.Context) {
		if ctx.Method() != fiber.MethodPost {
			ctx.Next()
			return
		}

		reqBody, err := ctx.GetBody()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}

		var encryptedBytes []byte
		if err = json.Unmarshal(reqBody, &encryptedBytes); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}

		decryptedBody, err := privateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}
		ctx.Request().Body = ioutil.NopCloser(bytes.NewReader(decryptedBody))
		ctx.Next()
	}
}

func SignResponse(privateKey *rsa.PrivateKey) iris.Handler {
	return func(ctx iris.Context) {
		rec := ctx.Recorder()
		ctx.Next()

		respBody := rec.Body()
		if len(respBody) == 0 {
			return
		}

		signatureBytes, err := apiutil.RsaSign(privateKey, respBody)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}

		ctx.Header(model.SignatureHeader, hex.EncodeToString(signatureBytes))
	}
}

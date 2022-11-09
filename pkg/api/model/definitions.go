package model

const (
	// SignatureHeader is the header that will contain an RSA signature of the response from the server we will use
	// to verify the response is legitimate.
	SignatureHeader = "X-Signature"

	// IgnoreCryptoHeader is the header that will tell the client to ignore crypto for the request
	IgnoreCryptoHeader = "X-Ignore-Crypto"

	// DecryptedBodyCtxKey is the context key under which the decrypted request body can be found in the server.
	DecryptedBodyCtxKey = "decryptedBody"

	// ErrorKey is the key in Context.Values where the error for a response can be found, if there is one.
	ErrorKey = "error"

	// WebsocketAuthType is the authentication type used in the Authorization header of the websocket upgrade request.
	WebsocketAuthType = "Nebula"
)

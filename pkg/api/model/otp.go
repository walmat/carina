package model

import "github.com/google/uuid"

type (
	OtpData struct {
		Secret        string    `json:"secret"`
		Uid           uuid.UUID `json:"uid"`
		RecoveryCodes []byte    `json:"recoveryCodes"`
	}

	OtpEnableRequest struct {
		SessionId uuid.UUID `json:"sessionId"`
	}

	OtpEnableResponse struct {
		Success       bool     `json:"success"`
		QrCode        string   `json:"qrCode,omitempty"`
		RecoveryCodes []string `json:"recoveryCodes,omitempty"`
	}

	OtpVerifyRequest struct {
		SessionId uuid.UUID `json:"sessionId"`
		Code      string    `json:"code"`
	}

	OtpVerifyResponse struct {
		Success bool `json:"success"`
	}
)

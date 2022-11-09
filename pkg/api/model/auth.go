package model

import (
	"github.com/google/uuid"
	"nebula/pkg/security"
)

type AuthRequest struct {
	Email string        `json:"email"`
	Pass  string        `json:"pass"`
	Hwid  security.HWID `json:"hwid"`
}

type AuthUser struct {
	Id           string `json:"id"`
	Hash         string `json:"hash"`
	Email        string `json:"email"`
	Type         string `json:"type"`
	Instances    int    `json:"instances"`
	MaxInstances int    `json:"maxInstances"`
}

type AuthResponse struct {
	Success bool `json:"success"`

	User      *AuthUser  `json:"user,omitempty"`
	SessionId *uuid.UUID `json:"sessionId,omitempty"`

	Message string `json:"message,omitempty"`
}

type ActivateSessionRequest struct {
	SessionId uuid.UUID     `json:"sid"`
	Hwid      security.HWID `json:"hwid"`
}

type DeactivateSessionRequest struct {
	SessionId uuid.UUID `json:"sid"`
}

type PasswordResetRequest struct {
	Email             string `json:"email"`
	RecaptchaResponse string `json:"recaptchaResponse"`
}

type PasswordResetResponse struct {
	Success bool `json:"success"`
}

type PasswordChangeRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

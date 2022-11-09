package model

import "github.com/google/uuid"

type RegisterRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	Key               string `json:"key"`
	RecaptchaResponse string `json:"recaptchaResponse"`
}

type RegisterResponse struct {
	Success bool `json:"success"`
}

type CompleteRegistrationRequest struct {
	Id uuid.UUID `json:"id"`
}

type CompleteRegistrationResponse struct {
	Success bool `json:"success"`
}

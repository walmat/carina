package rest

import (
	"encoding/json"
	"errors"
	"nebula/pkg/api/client"
	"nebula/pkg/api/model"
	"nebula/pkg/security"

	"github.com/google/uuid"
)

func AuthUser(email, pass string, hwid security.HWID) (*model.AuthResponse, error) {
	resp, err := client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(model.AuthRequest{
			Email: email,
			Pass:  pass,
			Hwid:  hwid,
		}).
		Post(getApiUrl() + "/api/auth")
	if err != nil {
		return nil, err
	}

	var authResponse model.AuthResponse
	if err = json.Unmarshal(resp.Body(), &authResponse); err != nil {
		return nil, err
	}

	if !authResponse.Success && authResponse.SessionId == nil {
		return nil, errors.New(authResponse.Message)
	}

	return &authResponse, nil
}

func ActivateSession(sid uuid.UUID, hwid security.HWID) (*model.AuthResponse, error) {
	resp, err := client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(model.ActivateSessionRequest{
			SessionId: sid,
			Hwid:      hwid,
		}).
		Post(getApiUrl() + "/api/auth/session")
	if err != nil {
		return nil, err
	}

	var sessionResponse model.AuthResponse
	if err = json.Unmarshal(resp.Body(), &sessionResponse); err != nil {
		return nil, err
	}

	if !sessionResponse.Success {
		return nil, errors.New(sessionResponse.Message)
	}

	return &sessionResponse, nil
}

func DeactivateSession(sid uuid.UUID) (*model.AuthResponse, error) {
	resp, err := client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(model.DeactivateSessionRequest{
			SessionId: sid,
		}).
		Post(getApiUrl() + "/api/auth/deactivate")
	if err != nil {
		return nil, err
	}

	var sessionResponse model.AuthResponse
	if err = json.Unmarshal(resp.Body(), &sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse, nil
}

func ForgotPassword(email, recaptchaResponse string) (*model.PasswordResetResponse, error) {
	resp, err := client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(model.PasswordResetRequest{
			Email:             email,
			RecaptchaResponse: recaptchaResponse,
		}).
		Post(getApiUrl() + "/api/auth/reset")
	if err != nil {
		return nil, err
	}

	var passwordResponse model.PasswordResetResponse
	if err = json.Unmarshal(resp.Body(), &passwordResponse); err != nil {
		return nil, err
	}

	return &passwordResponse, nil
}

func ChangePassword(token, password string) (*model.PasswordResetResponse, error) {
	resp, err := client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(model.PasswordChangeRequest{
			Token:    token,
			Password: password,
		}).
		Post(getApiUrl() + "/api/auth/change")
	if err != nil {
		return nil, err
	}

	var passwordResponse model.PasswordResetResponse
	if err = json.Unmarshal(resp.Body(), &passwordResponse); err != nil {
		return nil, err
	}

	return &passwordResponse, nil
}

func SubmitOtp(sid uuid.UUID, code string) (*model.OtpVerifyResponse, error) {
	resp, err := client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(model.OtpVerifyRequest{
			SessionId: sid,
			Code:      code,
		}).
		Post(getApiUrl() + "/api/auth/otp/verify")
	if err != nil {
		return nil, err
	}

	var otpResp model.OtpVerifyResponse
	if err = json.Unmarshal(resp.Body(), &otpResp); err != nil {
		return nil, err
	}

	return &otpResp, nil
}

func Register(key, email, password, recaptchaResponse string) (*model.RegisterResponse, error) {
	resp, err := client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(model.RegisterRequest{
			Email:             email,
			Password:          password,
			Key:               key,
			RecaptchaResponse: recaptchaResponse,
		}).
		Post(getApiUrl() + "/api/auth/register")
	if err != nil {
		return nil, err
	}

	var registerResponse model.RegisterResponse
	if err = json.Unmarshal(resp.Body(), &registerResponse); err != nil {
		return nil, err
	}

	return &registerResponse, nil
}

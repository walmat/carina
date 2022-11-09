package app

import (
	"fmt"
	"nebula/cmd/bot/windows"
	"nebula/pkg/api/client/rest"
	"nebula/pkg/api/client/ws"
	"nebula/pkg/clientutil"
	"nebula/pkg/security/hwid"
	"nebula/third_party/themida"
	"strings"
)

func authLogin(email, password string) LoginStatus {
	themida.Macro(themida.EAGLE_BLACK_START)

	// Login is used for logging in the user.
	authResp, err := rest.AuthUser(email, password, hwid.Get())
	//TODO: we should probably identify the issue further to clarify for the user.
	if err != nil {
		return LoginStatus{
			Message: strings.Title(err.Error()),
			Success: false,
		}
	} else if authResp.Message == "Submit OTP" {
		clientutil.SetSession(authResp.SessionId)

		port, err := GetUiPort(3001)
		if err != nil {
			return LoginStatus{
				Message: "Internal Error",
				Success: false,
			}
		}

		otpUrl := fmt.Sprintf("http://127.0.0.1:%d/login.html#2fa", port)

		if err = windows.Windows.LoginWindow.LoadURL(otpUrl); err != nil {
			return LoginStatus{
				Message: "Internal Error",
				Success: false,
			}
		}

		return LoginStatus{
			Success: false,
		}
	} else if !authResp.Success {
		return LoginStatus{
			Message: authResp.Message,
			Success: false,
		}
	}

	/*
		if err = ws.Connect(*authResp.SessionId); err != nil {
			return LoginStatus{
				Message: "Unknown Error",
				Success: false,
			}
		}
	*/
	clientutil.SetUser(authResp.User)
	clientutil.SetSession(authResp.SessionId)

	themida.Macro(themida.EAGLE_BLACK_END)

	TransitionToMain()

	return LoginStatus{
		Message: "Successfully Authenticated",
		Success: true,
	}
}

func submitOtp(code string) bool {
	sid := clientutil.GetSession()
	if sid == nil {
		return false
	}
	otpResp, err := rest.SubmitOtp(*sid, code)
	if err != nil || !otpResp.Success {
		return false
	}

	authResp, err := rest.ActivateSession(*sid, hwid.Get())
	if err != nil || !authResp.Success {
		return false
	}

	clientutil.SetUser(authResp.User)
	if err = ws.Connect(*sid); err != nil {
		return false
	}

	TransitionToMain()

	return otpResp.Success
}

func forgotPassword(email, recaptchaResponse string) (bool, error) {
	resp, err := rest.ForgotPassword(email, recaptchaResponse)
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

func changePassword(token, password string) (bool, error) {
	resp, err := rest.ChangePassword(token, password)
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

func register(key, email, password, recaptchaResponse string) LoginStatus {
	resp, err := rest.Register(key, email, password, recaptchaResponse)

	if err != nil {
		return LoginStatus{
			Success: false,
			Message: strings.Title(err.Error()),
		}
	}

	if !resp.Success {
		return LoginStatus{
			Success: false,
			Message: "There was an issue registering your account",
		}
	}

	return LoginStatus{
		Success: true,
		Message: "Email confirmation sent! Please check your inbox.",
	}
}

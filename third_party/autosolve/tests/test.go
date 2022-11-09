package main

import (
	"nebula/third_party/autosolve"
	"sync"
)

func main() {
	//Establishes listener function to receive token responses
	var tokenListener autosolve.CaptchaTokenResponseListener = func(tokenResponse autosolve.CaptchaTokenResponse) {
		print("Token response")
	}

	//Establishes listener function to receive cancel token responses
	var tokenCancelListener autosolve.CaptchaTokenCancelResponseListener = func(cancelTokenResponse autosolve.CaptchaTokenCancelResponse) {
		print("Response from cancel")
	}

	//Establishes listener function to receive status updates from AutoSolve
	var statusListener autosolve.StatusListener = func(status autosolve.Status) {

	}

	//Establishes listener function to receive errors
	var errorListener autosolve.ErrorListener = func(err error) {

	}

	err := autosolve.Load("", tokenListener, tokenCancelListener, statusListener, errorListener)

	if err != nil {
		print("Error in Load")
	}

	res, err := autosolve.Connect("", "")

	switch res {
	case autosolve.Success:
		print("Successful Connection")
		var message = autosolve.CaptchaTokenRequest{
			TaskId:  "1",
			Url:     "https://recaptcha.autosolve.io/version/1",
			SiteKey: "6Ld_LMAUAAAAAOIqLSy5XY9-DUKLkAgiDpqtTJ9b",
		}
		autosolve.SendTokenRequest(message)
		message.TaskId = "7"
		autosolve.SendTokenRequest(message)
		wg := sync.WaitGroup{}
		wg.Add(1)
		wg.Wait()
	case autosolve.InvalidClientId:
		print("Invalid Client Key")
	case autosolve.InvalidAccessToken:
		print("Invalid access token")
	case autosolve.InvalidApiKey:
		print("Invalid Api Key")
	case autosolve.InvalidCredentials:
		print("Invalid Credentials")
	}
}

package autosolve

import (
	"errors"
	"nebula/pkg/integrations"
	"nebula/pkg/logger"
	"nebula/third_party/autosolve"
	"sync"
	"time"
)

var (
	callbacks     = make(map[string]chan CaptchaToken)
	callbackMutex = sync.RWMutex{}
)

type CaptchaToken autosolve.CaptchaTokenResponse

func (t *CaptchaToken) IsExpired() bool {
	switch t.Request.Version {
	case CaptchaTypeRecap3:
	case CaptchaTypeRecap2Checkbox:
	case CaptchaTypeRecap2Invisible:
		const TwoMinutes = int64(time.Minute*2) / 1000
		return time.Now().Unix()-t.CreatedAt > TwoMinutes
	}

	return false
}

type CaptchaRequest autosolve.CaptchaTokenRequest

type CaptchaType int

const (
	CaptchaTypeRecap2Checkbox  = 0
	CaptchaTypeRecap2Invisible = 1
	CaptchaTypeRecap3          = 2
	CaptchaTypeHCapCheckbox    = 3
	CaptchaTypeHCapInvisible   = 4
	CaptchaTypeGeeTest         = 5
)

func init() {
	tokenListener := func(tokenResponse autosolve.CaptchaTokenResponse) {
		callbackMutex.Lock()
		defer callbackMutex.Unlock()
		if callback, ok := callbacks[tokenResponse.TaskId]; ok {
			callback <- CaptchaToken(tokenResponse)
			delete(callbacks, tokenResponse.TaskId)
		}
	}

	tokenCancelListener := func(cancelTokenResponse autosolve.CaptchaTokenCancelResponse) {
		callbackMutex.Lock()
		defer callbackMutex.Unlock()
		for _, request := range cancelTokenResponse.Requests {
			if callback, ok := callbacks[request.TaskId]; ok {
				close(callback)
				delete(callbacks, request.TaskId)
			}
		}
	}

	statusListener := func(status autosolve.Status) {
		// TODO: do we need to handle this?
		// NOTE: Sends Connecting | Connected
	}

	errorListener := func(err error) {
		logger.Error("autosolve error:", err)
	}

	if err := autosolve.Load("Nebula-c92504a1-5441-4970-9218-be520bc5416c", tokenListener, tokenCancelListener, statusListener, errorListener); err != nil {
		// TODO: should this crash app as well?
		logger.Error("failed to initialize autosolve:", err)
	}
}

func Connect(accessToken, apiKey string) error {
	result, err := autosolve.Connect(accessToken, apiKey)
	if err != nil {
		return err
	}

	switch result {
	case autosolve.InvalidClientId:
		return errors.New("invalid client id")
	case autosolve.InvalidAccessToken:
		return errors.New("invalid access token")
	case autosolve.InvalidApiKey:
		return errors.New("invalid api key")
	case autosolve.InvalidCredentials:
		return errors.New("invalid credentials")
	default: // basically autosolve.Success
		integrations.SetAutoSolve(&[]integrations.Credential{
			{
				Label: "Access Token",
				Value: accessToken,
			},
			{
				Label: "Api Key",
				Value: apiKey,
			},
		})

		return nil
	}
}

func Close() error {
	integrations.SetAutoSolve(nil)
	return autosolve.Close()
}

func RequestCaptchaToken(req CaptchaRequest) (chan CaptchaToken, error) {
	callback := make(chan CaptchaToken)
	callbackMutex.Lock()
	callbacks[req.TaskId] = callback
	callbackMutex.Unlock()
	return callback, autosolve.SendTokenRequest(autosolve.CaptchaTokenRequest(req))
}

func CancelCaptchaRequest(taskId string) error {
	callbackMutex.Lock()
	defer callbackMutex.Unlock()
	if callback, ok := callbacks[taskId]; ok {
		close(callback)
		delete(callbacks, taskId)
	}
	return autosolve.SendTokenCancelRequest(autosolve.CaptchaTokenCancelRequest{
		TaskId:           taskId,
		ResponseRequired: false,
	})
}

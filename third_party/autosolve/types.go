package autosolve

type CaptchaTokenRequest struct {
	TaskId           string            `json:"taskId"`
	ApiKey           string            `json:"apiKey"`
	CreatedAt        int64             `json:"createdAt"`
	Url              string            `json:"url"`
	SiteKey          string            `json:"siteKey"`
	Version          int               `json:"version"`
	Action           string            `json:"action"`
	MinScore         float32           `json:"minScore"`
	Proxy            string            `json:"proxy"`
	ProxyRequired    bool              `json:"proxyRequired"`
	UserAgent        string            `json:"userAgent"`
	Cookies          string            `json:"cookies"`
	RenderParameters map[string]string `json:"renderParameters"`
}

type CaptchaTokenCancelRequest struct {
	TaskId           string   `json:"taskId"`
	ApiKey           string   `json:"apiKey"`
	CreatedAt        int64    `json:"createdAt"`
	TaskIds          []string `json:"taskIds"`
	ResponseRequired bool     `json:"responseRequired"`
}

type CaptchaTokenResponse struct {
	TaskId    string              `json:"taskId"`
	ApiKey    string              `json:"apiKey"`
	CreatedAt int64               `json:"createdAt"`
	Request   CaptchaTokenRequest `json:"request"`
	Token     string              `json:"token"`
}

type CaptchaTokenCancelResponse struct {
	Requests []CaptchaTokenRequest `json:"requests"`
}

type Status string

const (
	Connecting   Status = "Connecting"
	Connected    Status = "Connected"
	Reconnecting Status = "Reconnecting"
	Disconnected Status = "Disconnected"
)

type CaptchaTokenResponseListener func(response CaptchaTokenResponse)
type CaptchaTokenCancelResponseListener func(response CaptchaTokenCancelResponse)
type ErrorListener func(err error)
type StatusListener func(status Status)

type ConnectResult string

const (
	Success            ConnectResult = "Success"
	InvalidAccessToken ConnectResult = "InvalidAccessToken"
	InvalidApiKey      ConnectResult = "InvalidApiKey"
	InvalidCredentials ConnectResult = "InvalidCredentials"
	InvalidClientId    ConnectResult = "InvalidClientId"
	ConnectionPending  ConnectResult = "ConnectionPending"
	TooManyRequests    ConnectResult = "TooManyRequests"
	ConnectionError    ConnectResult = "ConnectionFailed"
	VerificationError  ConnectResult = "VerificationFailed"
	UnknownError       ConnectResult = "UnknownError"
)

type Account struct {
	id           int
	rId          string
	accessToken  string
	rAccessToken string
	apiKey       string
	rApiKey      string
}

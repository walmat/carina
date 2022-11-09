package harvesters

import (
	"encoding/json"
	"nebula/pkg/integrations/autosolve"
	"nebula/third_party/http"
	"net/url"

	"github.com/Nebulabots/go-astilectron"
)

func UnmarshalHarvesterData(data []byte) (HarvesterData, error) {
	var harvester HarvesterData
	err := json.Unmarshal(data, &harvester)
	return harvester, err
}

type Store struct {
	Id       string `json:"id"`
	Platform string `json:"platform"`
	Name     string `json:"name"`
	Url      string `json:"url"`
}

type Account struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Parameters struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ReCaptcha struct {
	Parameters map[string]string `json:"parameters,omitempty"`
	SiteKey    string            `json:"siteKey"`
	Action     string            `json:"action,omitempty"`
}

type Solver struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"`
	Store       *Store         `json:"store"`
	ReCaptcha   *ReCaptcha     `json:"reCaptcha.omitempty"`
	Cookies     []*http.Cookie `json:"cookies,omitempty"`
	Proxy       *string        `json:"proxy,omitempty"`
	RequestedAt int64          `json:"requestedAt"`
}

type HarvesterData struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Index   int        `json:"index"`
	Focused bool       `json:"focused"`
	Store   *Store     `json:"store"`
	Account *GmailData `json:"account"`
	Proxy   string     `json:"proxy"`
	Solver  *Solver    `json:"solver"`
}

type ResponseData struct {
	Token   string                      `json:"token,omitempty"`
	Cookies []astilectron.SessionCookie `json:"cookies,omitempty"`
	Email   string                      `json:"email,omitempty"`
	Form    url.Values                  `json:"form,omitempty"`
}

type SolveChans struct {
	Autosolve  chan autosolve.CaptchaToken
	Collective chan SolverResponse
}

type SolverResponse struct {
	ID        string       `json:"id"`
	Success   bool         `json:"success"`
	CreatedAt int64        `json:"createdAt"`
	Data      ResponseData `json:"data"`
}

type GmailData struct {
	ID      string                       `json:"id"`
	Email   string                       `json:"email"`
	Cookies *[]astilectron.SessionCookie `json:"cookies"`
}

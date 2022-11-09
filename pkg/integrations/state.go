package integrations

import (
	"encoding/json"
	"os"
)

var (
	integrations = Integrations{
		Aycd: Integration{
			ID:     "aycd",
			Name:   "AYCD AutoSolve",
			Active: false,
			Credentials: []Credential{
				{
					Label: "Access Token",
					Value: "",
				},
				{
					Label: "Api Key",
					Value: "",
				},
			},
		},
		TwoCaptcha: Integration{
			ID:     "2captcha",
			Name:   "2Captcha",
			Active: false,
			Credentials: []Credential{
				{
					Label: "Api Key",
					Value: "",
				},
			},
		},
		CapMonster: Integration{
			ID:     "capmonster",
			Name:   "Cap Monster",
			Active: false,
			Credentials: []Credential{
				{
					Label: "Api Key",
					Value: "",
				},
			},
		},
		Scout: Integration{
			ID:          "scout",
			Name:        "Scout Analytics",
			Active:      false,
			Credentials: nil,
		},
	}
)

type Credential struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Integrations struct {
	Aycd       Integration `json:"aycd"`
	TwoCaptcha Integration `json:"2captcha"`
	CapMonster Integration `json:"capmonster"`
	Scout      Integration `json:"scout"`
}

type Integration struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Active      bool         `json:"active"`
	Credentials []Credential `json:"credentials"`
}

func GetAutoSolve() Integration {
	return integrations.Aycd
}

func SetAutoSolve(creds *[]Credential) {
	if creds != nil {
		integrations.Aycd.Active = true
		integrations.Aycd.Credentials = *creds
	} else {
		integrations.Aycd.Active = false
		integrations.Aycd.Credentials = []Credential{
			{
				Label: "Access Token",
				Value: "",
			},
			{
				Label: "Api Key",
				Value: "",
			},
		}
	}
}

func GetTwoCaptcha() Integration {
	return integrations.TwoCaptcha
}

func GetCapMonster() Integration {
	return integrations.CapMonster
}

func GetScout() Integration {
	return integrations.Scout
}

func Get() Integrations {
	return integrations
}

func Set(i Integrations) {
	integrations = i
}

func ImportIntegrations(outPath string) error {
	taskBytes, err := os.ReadFile(outPath)
	if err != nil {
		return err
	}

	var i Integrations
	if err = json.Unmarshal(taskBytes, &i); err != nil {
		return err
	}

	integrations = i
	return nil
}

func ExportIntegrations(outPath string) error {
	exportedBytes, err := json.Marshal(integrations)
	if err != nil {
		return err
	}

	return os.WriteFile(outPath, exportedBytes, os.ModePerm)
}

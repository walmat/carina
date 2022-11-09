package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"nebula/pkg/api/client"
	"runtime"
)

func Download() ([]byte, error) {
	var userAgent string
	switch runtime.GOOS {
	case "windows":
		userAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"
	case "darwin":
		userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8"
	default:
		return nil, errors.New("improper runtime")
	}

	resp, err := client.
		R().
		SetHeader("user-agent", userAgent).
		Get(getApiUrl() + "/api/download")
	if err != nil {
		return nil, err
	}

	if err != nil {
		return resp.Body(), err
	}

	return resp.Body(), nil
}

type updateResponse struct {
	Latest string `json:"latest"`
}

func CheckForUpdates() (updateResponse, error) {

	resp, err := client.
		R().
		Get(getApiUrl() + fmt.Sprintf("/api/update?platform=%s", runtime.GOOS))
	if err != nil {
		return updateResponse{}, err
	}

	var res updateResponse
	err = json.Unmarshal(resp.Body(), &res)

	if err != nil {
		return updateResponse{}, err
	}

	return res, nil
}

package rest

import (
	"fmt"
	"nebula/pkg/api/client"
	"nebula/pkg/api/model"
	"nebula/pkg/clientutil"
)

func Log(entry model.LoggerEntry) error {
	_, err := client.
		R().
		SetHeader("Authorization", fmt.Sprintf("Nebula %s", clientutil.GetSession().String())).
		SetHeader("Content-Type", "application/json").
		SetBody(entry).
		Post(getApiUrl() + "/api/log")
	return err
}

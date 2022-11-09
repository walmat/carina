package antibot

import (
	"encoding/json"
	"errors"
	"fmt"
	"nebula/third_party/http"
	"strconv"
	"strings"
)

type PerimeterXGenerator struct {
	UserAgent string

	payload string

	uuid string

	PXCookie string

	timestamp int

	response map[string]interface{}
}

func (px *PerimeterXGenerator) Generate(client *http.SimpleClient) error {
	px.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"

	if err := px.GetFirstPayload(client); err != nil {
		return fmt.Errorf("error: get first payload: %s", err)
	} else if err = px.PostFirstPayload(client); err != nil {
		return fmt.Errorf("error: post first payload: %s", err)
	} else if err = px.GetSecondPayload(client); err != nil {
		return fmt.Errorf("error: get second payload: %s", err)
	} else if err = px.PostSecondPayload(client); err != nil {
		return fmt.Errorf("error: post second payload: %s", err)
	}

	return nil
}

func (px *PerimeterXGenerator) GetPXFirstPayload(client *http.SimpleClient) (map[string]interface{}, error) {
	resp, err := client.Builder("https://nebula-test-perimeterx-server-twdug.ondigitalocean.app/post_first").
		Method(http.MethodPost).
		Header("user-agent", px.UserAgent).
		JsonBody(map[string]string{
			"domain":    "https://www.walmart.com/",
			"userAgent": px.UserAgent,
			"appId":     "PXu6b0qd2S",
		}).
		SendAndClose()

	if err != nil {
		return nil, err
	} else if resp.StatusCode > 201 {
		return nil, errors.New("invalid status code")
	}

	jsonData := map[string]interface{}{}
	if err = json.Unmarshal(resp.BytesBody(), &jsonData); err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (px *PerimeterXGenerator) GetPXSecondPayload(client *http.SimpleClient) (map[string]interface{}, error) {
	payloadMarshalled, err := json.Marshal(px.response)
	if err != nil {
		return nil, err
	}

	resp, err := client.Builder("https://nebula-test-perimeterx-server-twdug.ondigitalocean.app/post_first").
		Method(http.MethodPost).
		Header("user-agent", px.UserAgent).
		JsonBody(map[string]string{
			"domain":    "https://www.walmart.com/",
			"userAgent": px.UserAgent,
			"appId":     "PXu6b0qd2S",
			"payload":   string(payloadMarshalled),
			"uuid":      px.uuid,
			"timestamp": strconv.Itoa(px.timestamp),
		}).
		SendAndClose()

	if err != nil {
		return nil, err
	} else if resp.StatusCode > 201 {
		return map[string]interface{}{}, errors.New("invalid status code")
	}

	jsonData := map[string]interface{}{}
	err = json.Unmarshal(resp.BytesBody(), &jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (px *PerimeterXGenerator) GetFirstPayload(client *http.SimpleClient) error {
	beaconData, err := px.GetPXFirstPayload(client)
	if err != nil {
		return err
	}

	results, resultsOk := beaconData["result"]

	if !resultsOk || results == nil {
		return errors.New("error generating perimeterx")
	}

	px.payload = ""
	var arrayValues []string
	for key, value := range results.(map[string]interface{}) {
		arrayValues = append(arrayValues, key+"="+value.(string))
	}

	px.payload = strings.Join(arrayValues, "&")

	timestamp, ok := beaconData["ts"]
	if !ok || timestamp == nil {
		return errors.New("error generating perimeterx")
	}

	timestampFloat, ok := timestamp.(float64)
	if !ok {
		return errors.New("error generating perimeterx")
	}

	px.timestamp = int(timestampFloat)
	px.uuid = results.(map[string]interface{})["uuid"].(string)

	return nil
}

func (px *PerimeterXGenerator) PostFirstPayload(client *http.SimpleClient) error {
	resp, err := client.Builder("https://nebula-test-perimeterx-server-twdug.ondigitalocean.app/post_first").
		Method(http.MethodPost).
		Header("user-agent", px.UserAgent).
		Header("accept", "*/*").
		Header("origin", "https://www.walmart.com").
		Header("referer", "https://www.walmart.com/").
		Header("sec-ch-ua", `"Google Chrome";v="89", "Chromium";v="89", ";Not A Brand";v="99"`).
		Header("sec-ch-ua-mobile", "?0").
		Header("sec-fetch-site", "cross-site").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		Header("content-type", "application/x-www-form-urlencoded").
		Header("content-length", strconv.Itoa(len(px.payload))).
		HeaderOrder("accept", "accept-encoding", "accept-language", "content-length", "content-type", "origin", "referer", "sec-ch-ua", "sec-ch-ua-mobile", "sec-fetch-dest", "sec-fetch-mode", "sec-fetch-site", "user-agent").
		StringBody(px.payload).
		Send()
	if err != nil {
		return err
	} else if resp.StatusCode > 201 {
		return errors.New("invalid status code")
	}

	err = json.Unmarshal(resp.BytesBody(), &px.response)
	if err != nil {
		return err
	}

	return nil
}

func (px *PerimeterXGenerator) GetSecondPayload(client *http.SimpleClient) error {
	beaconData, beaconErr := px.GetPXSecondPayload(client)

	if beaconErr != nil {
		return beaconErr
	}

	px.payload = ""
	var arrayValues []string
	for key, value := range beaconData {
		arrayValues = append(arrayValues, key+"="+value.(string))
	}

	px.payload = strings.Join(arrayValues, "&")

	return nil
}

func (px *PerimeterXGenerator) PostSecondPayload(client *http.SimpleClient) error {
	resp, err := client.Builder("https://nebula-test-perimeterx-server-twdug.ondigitalocean.app/post_first").
		Method(http.MethodPost).
		Header("user-agent", px.UserAgent).
		Header("accept", "*/*").
		Header("origin", "https://www.walmart.com").
		Header("referer", "https://www.walmart.com/").
		Header("sec-ch-ua", `"Google Chrome";v="89", "Chromium";v="89", ";Not A Brand";v="99"`).
		Header("sec-ch-ua-mobile", "?0").
		Header("sec-fetch-site", "cross-site").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		Header("content-type", "application/x-www-form-urlencoded").
		Header("content-length", strconv.Itoa(len(px.payload))).
		HeaderOrder("accept", "accept-encoding", "accept-language", "content-length", "content-type", "origin", "referer", "sec-ch-ua", "sec-ch-ua-mobile", "sec-fetch-dest", "sec-fetch-mode", "sec-fetch-site", "user-agent").
		StringBody(px.payload).
		SendAndClose()

	if err != nil {
		return err
	}

	if resp.StatusCode > 201 {
		return errors.New("invalid status code")
	}

	responseBody := make(map[string]interface{})
	if err = json.Unmarshal(resp.BytesBody(), &responseBody); err != nil {
		return err
	}

	do, doExists := responseBody["do"]
	if !doExists || do == nil {
		return errors.New("do does not exist")
	}

	doArray := do.([]interface{})
	if len(doArray) == 0 {
		return errors.New("do array empty")
	}

	for _, doValue := range doArray {
		doString := doValue.(string)

		parts := strings.Split(doString, "|")

		if parts[0] == "bake" {
			px.PXCookie = parts[3]
		}
	}

	if len(px.PXCookie) == 0 {
		return errors.New("no px cookie generated")
	}

	return nil
}

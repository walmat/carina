package rest

import "nebula/pkg/api/client"

func getApiUrl() string {
	apiClientUrl := client.ApiUrl
	if apiClientUrl == "" {
		apiClientUrl = "http://localhost:8080"
	}
	return apiClientUrl
}

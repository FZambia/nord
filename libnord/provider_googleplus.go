package libnord

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type GoogleplusProvider struct {
	client *http.Client
}

func NewGoogleplusProvider(client *http.Client) *GoogleplusProvider {
	return &GoogleplusProvider{
		client: client,
	}
}

func (provider *GoogleplusProvider) GetUrlInfo(url string) *ProviderResponse {

	resp := &ProviderResponse{
		Name: "googleplus",
		Url:  url,
	}

	rpc_objects := make([]map[string]interface{}, 0)

	params := map[string]interface{}{
		"nolog":   true,
		"id":      url,
		"source":  "widget",
		"userId":  "@viewer",
		"groupId": "@self",
	}

	rpc_obj := map[string]interface{}{
		"method":     "pos.plusones.get",
		"id":         "go",
		"jsonrpc":    "2.0",
		"key":        "go",
		"apiVersion": "v1",
		"params":     params,
	}

	rpc_objects = append(rpc_objects, rpc_obj)

	baseUrl := "https://clients6.google.com/rpc"
	jsonBody, err := json.Marshal(rpc_objects)

	bodyBuffer := bytes.NewBuffer(jsonBody)

	content, err := makePostRequest(provider.client, baseUrl, "application/json", bodyBuffer)
	if err != nil {
		resp.Err = err.Error()
		return resp
	}

	var f interface{}
	err = json.Unmarshal(content, &f)
	m := f.([]interface{})
	if len(m) == 0 {
		resp.Err = ERROR_NOT_FOUND
	}
	obj := m[0]
	dict := obj.(map[string]interface{})
	_, exists := dict["error"]
	if exists {
		resp.Err = ERROR_NOT_FOUND
		return resp
	}
	result := dict["result"].(map[string]interface{})
	metadata := result["metadata"].(map[string]interface{})
	counts := metadata["globalCounts"].(map[string]interface{})
	resp.Data = map[string]int{
		"count": int(counts["count"].(float64)),
	}

	return resp
}

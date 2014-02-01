package libnord

import (
	"encoding/json"
	"net/http"
)

type PinterestProvider struct {
	client *http.Client
}

func NewPinterestProvider(client *http.Client) *PinterestProvider {
	return &PinterestProvider{
		client: client,
	}
}

func (provider *PinterestProvider) GetUrlInfo(url string) *ProviderResponse {

	params := map[string]string{
		"callback": "go",
		"url":      url,
	}

	baseUrl := "https://api.pinterest.com/v1/urls/count.json"
	endpoint := urlencode(baseUrl, params)

	resp := &ProviderResponse{
		Name: "pinterest",
		Url:  url,
	}

	content, err := makeGetRequest(provider.client, endpoint)
	if err != nil {
		resp.Err = ERROR_NOT_FOUND
		return resp
	}

	// unwrap jsonp callback
	content = content[3 : len(content)-1]

	var f interface{}
	err = json.Unmarshal(content, &f)
	m := f.(map[string]interface{})
	_, exists := m["error"]
	if exists {
		resp.Err = ERROR_NOT_FOUND
		return resp
	}
	resp.Data = map[string]int{
		"count": int(m["count"].(float64)),
	}
	return resp
}

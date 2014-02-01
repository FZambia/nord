package libnord

import (
	"encoding/json"
	"net/http"
)

type TwitterProvider struct {
	client *http.Client
}

func NewTwitterProvider(client *http.Client) *TwitterProvider {
	return &TwitterProvider{
		client: client,
	}
}

func (provider *TwitterProvider) GetUrlInfo(url string) *ProviderResponse {

	params := map[string]string{
		"url": url,
	}

	baseUrl := "http://urls.api.twitter.com/1/urls/count.json"
	endpoint := urlencode(baseUrl, params)

	resp := &ProviderResponse{
		Name: "twitter",
		Url:  url,
	}

	content, err := makeGetRequest(provider.client, endpoint)
	if err != nil {
		resp.Err = ERROR_NOT_FOUND
		return resp
	}

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

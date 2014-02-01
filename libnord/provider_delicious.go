package libnord

import (
	"encoding/json"
	"net/http"
)

type DeliciousProvider struct {
	client *http.Client
}

func NewDeliciousProvider(client *http.Client) *DeliciousProvider {
	return &DeliciousProvider{
		client: client,
	}
}

func (provider *DeliciousProvider) GetUrlInfo(url string) *ProviderResponse {

	params := map[string]string{
		"url": url,
	}

	baseUrl := "http://feeds.delicious.com/v2/json/urlinfo/data"
	endpoint := urlencode(baseUrl, params)

	resp := &ProviderResponse{
		Name: "delicious",
		Url:  url,
	}

	content, err := makeGetRequest(provider.client, endpoint)
	if err != nil {
		resp.Err = ERROR_NOT_FOUND
		return resp
	}

	var f interface{}
	err = json.Unmarshal(content, &f)
	m := f.([]interface{})

	if len(m) == 0 {
		resp.Err = "not found"
		return resp
	}
	obj := m[0].(map[string]interface{})

	_, exists := obj["error"]
	if exists {
		resp.Err = ERROR_NOT_FOUND
		return resp
	}
	resp.Data = map[string]int{
		"bookmarks": int(obj["total_posts"].(float64)),
	}
	return resp
}

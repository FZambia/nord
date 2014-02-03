package libnord

import (
	"encoding/json"
	"net/http"
)

type LinkedinProvider struct {
	client *http.Client
}

func NewLinkedinProvider(client *http.Client) *LinkedinProvider {
	return &LinkedinProvider{
		client: client,
	}
}

func (provider *LinkedinProvider) GetUrlInfo(url string) *ProviderResponse {

	params := map[string]string{
		"callback": "go",
		"url":      url,
	}

	baseUrl := "http://www.linkedin.com/countserv/count/share"
	endpoint := urlencode(baseUrl, params)

	resp := &ProviderResponse{
		Name: "linkedin",
		Url:  url,
	}

	content, err := makeGetRequest(provider.client, endpoint)
	if err != nil {
		resp.Err = ERROR_NOT_FOUND
		return resp
	}

	// unwrap jsonp callback
	content = content[3 : len(content)-2]

	var f interface{}
	err = json.Unmarshal(content, &f)
	m := f.(map[string]interface{})
	resp.Data = map[string]int{
		"share_count": int(m["count"].(float64)),
	}
	return resp
}

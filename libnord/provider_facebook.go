package libnord

import (
	"encoding/json"
	"net/http"
)

type FacebookProvider struct {
	client *http.Client
}

func NewFacebookProvider(client *http.Client) *FacebookProvider {
	return &FacebookProvider{
		client: client,
	}
}

func (provider *FacebookProvider) GetUrlInfo(url string) *ProviderResponse {

	params := map[string]string{
		"method": "links.getStats",
		"format": "json",
		"urls":   url,
	}

	baseUrl := "http://api.facebook.com/restserver.php"
	endpoint := urlencode(baseUrl, params)

	resp := &ProviderResponse{
		Name: "facebook",
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
	obj := m[0]

	dict := obj.(map[string]interface{})
	resp.Data = map[string]int{
		"like_count":        int(dict["like_count"].(float64)),
		"share_count":       int(dict["share_count"].(float64)),
		"comment_count":     int(dict["comment_count"].(float64)),
		"total_count":       int(dict["total_count"].(float64)),
		"commentsbox_count": int(dict["commentsbox_count"].(float64)),
		"click_count":       int(dict["click_count"].(float64)),
	}
	return resp
}

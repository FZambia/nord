package libnord

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type OkProvider struct {
	client *http.Client
}

func NewOkProvider(client *http.Client) *OkProvider {
	return &OkProvider{
		client: client,
	}
}

func (provider *OkProvider) GetUrlInfo(url string) *ProviderResponse {

	params := map[string]string{
		"cb":     "go",
		"st.cmd": "shareData",
		"ref":    url,
	}

	baseUrl := "https://odnoklassniki.ru/dk"
	endpoint := urlencode(baseUrl, params)

	resp := &ProviderResponse{
		Name: "ok",
		Url:  url,
	}

	content, err := makeGetRequest(provider.client, endpoint)
	if err != nil {
		resp.Err = ERROR_NOT_FOUND
		return resp
	}

	content = content[3 : len(content)-1]

	var f interface{}
	err = json.Unmarshal(content, &f)
	m := f.(map[string]interface{})
	count, err := strconv.Atoi(m["count"].(string))
	if err != nil {
		resp.Err = ERROR_BAD_RESPONSE
		return resp
	}
	resp.Data = map[string]int{
		"count": count,
	}
	return resp
}

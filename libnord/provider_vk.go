package libnord

import (
	"net/http"
	"regexp"
	"strconv"
)

type VkProvider struct {
	client *http.Client
}

func NewVkProvider(client *http.Client) *VkProvider {
	return &VkProvider{
		client: client,
	}
}

func (provider *VkProvider) GetUrlInfo(url string) *ProviderResponse {

	params := map[string]string{
		"act":   "count",
		"url":   url,
		"index": "1",
	}

	baseUrl := "https://vk.com/share.php"
	endpoint := urlencode(baseUrl, params)

	resp := &ProviderResponse{
		Name: "vk",
		Url:  url,
	}

	content, err := makeGetRequest(provider.client, endpoint)
	if err != nil {
		resp.Err = ERROR_NOT_FOUND
		return resp
	}

	re := regexp.MustCompile("VK.Share.count\\(1, ([0-9]+)\\);")
	matches := re.FindStringSubmatch(string(content))
	if matches == nil {
		resp.Err = ERROR_NOT_FOUND
		return resp
	}
	count, _ := strconv.Atoi(matches[1])
	resp.Data = map[string]int{
		"count": count,
	}
	return resp
}

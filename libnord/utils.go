package libnord

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func urlencode(baseUrl string, arguments map[string]string) string {
	params := url.Values{}
	for k, v := range arguments {
		params.Add(k, v)
	}
	endpoint := baseUrl + "?" + params.Encode()
	return endpoint
}

func makeGetRequest(client *http.Client, endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Close = true

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func makePostRequest(client *http.Client, endpoint string, bodyType string, body io.Reader) ([]byte, error) {
	response, err := http.Post(endpoint, bodyType, body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

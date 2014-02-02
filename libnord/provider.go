package libnord

import (
	"errors"
	"net/http"
)

const (
	ERROR_NOT_FOUND     = "not found"
	ERROR_TIMED_OUT     = "timed out"
	ERROR_NOT_SUPPORTED = "not supported"
	ERROR_BAD_RESPONSE  = "bad response"
)

type ProviderResponse struct {
	Name string      `json:"-"`
	Url  string      `json:"-"`
	Data interface{} `json:"data"`
	Err  interface{} `json:"error"`
}

type Provider interface {
	GetUrlInfo(url string) *ProviderResponse
}

func GetProviderByName(name string, client *http.Client) (Provider, error) {
	switch name {
	case "facebook":
		return NewFacebookProvider(client), nil
	case "googleplus":
		return NewGoogleplusProvider(client), nil
	case "twitter":
		return NewTwitterProvider(client), nil
	case "pinterest":
		return NewPinterestProvider(client), nil
	case "delicious":
		return NewDeliciousProvider(client), nil
	case "linkedin":
		return NewLinkedinProvider(client), nil
	case "vk":
		return NewVkProvider(client), nil
	case "ok":
		return NewOkProvider(client), nil
	default:
		return nil, errors.New(ERROR_NOT_SUPPORTED)
	}
}

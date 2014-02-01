package libnord

import (
	"errors"
	"net/http"
)

const (
	ERROR_NOT_FOUND     = "not found"
	ERROR_TIMED_OUT     = "timed out"
	ERROR_NOT_SUPPORTED = "not supported"
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
	if name == "facebook" {
		return NewFacebookProvider(client), nil
	} else if name == "googleplus" {
		return NewGoogleplusProvider(client), nil
	} else if name == "twitter" {
		return NewTwitterProvider(client), nil
	} else if name == "pinterest" {
		return NewPinterestProvider(client), nil
	} else if name == "delicious" {
		return NewDeliciousProvider(client), nil
	} else if name == "linkedin" {
		return NewLinkedinProvider(client), nil
	} else if name == "vk" {
		return NewVkProvider(client), nil
	} else {
		return nil, errors.New(ERROR_NOT_SUPPORTED)
	}
}

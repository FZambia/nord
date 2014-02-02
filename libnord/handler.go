package libnord

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/martini"
	"github.com/garyburd/redigo/redis"
)

type ServiceMap struct {
	client *http.Client
	config *Config
	conn   redis.Conn
}

func getRequestProviders(r *http.Request) ([]string, error) {
	providersArgs, ok := r.Form["providers"]
	if !ok {
		return nil, errors.New("providers required")
	}
	provider_list := strings.Split(providersArgs[0], ",")
	return provider_list, nil
}

func getRequestUrl(r *http.Request) (string, error) {
	urlArgs, ok := r.Form["url"]
	if !ok {
		return "", errors.New("url required")
	}
	url := urlArgs[0]
	return url, nil
}

func getRequestTimeout(r *http.Request, conf *Config) (time.Duration, error) {
	var timeout time.Duration
	timeoutArgs, ok := r.Form["timeout"]
	if ok {
		timeoutInt, err := strconv.Atoi(timeoutArgs[0])
		if err != nil {
			return time.Duration(0), errors.New("wrong timeout")
		}
		timeout = time.Duration(timeoutInt) * time.Millisecond
	} else {
		timeout = conf.Timeout
	}
	return timeout, nil
}

func getRequestCallback(r *http.Request) (string, error) {
	var callback string
	callbackArgs, ok := r.Form["callback"]
	if ok {
		callback = callbackArgs[0]
		match, _ := regexp.MatchString("^[A-z0-9_]+$", callback)
		if !match {
			return "", errors.New("wrong callback")
		}
	}
	return callback, nil
}

func isResponseReady(response map[string]interface{}, providerList []string) bool {
	return len(response) == len(providerList)
}

func NordHandler(service *ServiceMap, w http.ResponseWriter, r *http.Request) (int, string) {

	// extract request parameters from request
	r.ParseForm()
	providerList, err := getRequestProviders(r)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}
	url, err := getRequestUrl(r)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}
	timeout, err := getRequestTimeout(r, service.config)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}
	callback, err := getRequestCallback(r)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}

	// initialize response: key - provider name, value - provider response
	response := map[string]interface{}{}

	// create channel to receive provider's responses from goroutines
	result_channel := make(chan *ProviderResponse)

	for _, providerName := range providerList {
		provider, err := GetProviderByName(providerName, service.client)
		if err != nil {
			response[providerName] = &ProviderResponse{
				Err: err.Error(),
			}
			continue
		}

		// try to find response in cache
		cached_result, err := getCacheResult(service, providerName, url)
		if err == nil && cached_result != nil {
			// success - we found cached result and can use it
			response[providerName] = cached_result
			continue
		}

		// create actual request for provider information in separate goroutine
		go func() {
			resp := provider.GetUrlInfo(url)
			result_channel <- resp
		}()
	}

	// collect responses from goroutines
	if !isResponseReady(response, providerList) {
		quit := false
		for {
			select {
			case res := <-result_channel:
				response[res.Name] = res
				setCacheResult(service, res.Name, url, res)
				if isResponseReady(response, providerList) {
					quit = true
				}
			case <-time.After(timeout):
				quit = true
			}
			if quit {
				break
			}
		}
	}

	// extend response with timeouted data
	if !isResponseReady(response, providerList) {
		for _, providerName := range providerList {
			_, exists := response[providerName]
			if exists {
				continue
			}
			response[providerName] = &ProviderResponse{
				Name: providerName,
				Err:  ERROR_TIMED_OUT,
			}
		}
	}

	// prepare final response body
	result_json, _ := json.Marshal(response)
	result := string(result_json)
	if len(callback) > 0 {
		w.Header().Set("Content-Type", "application/javascript")
		result = fmt.Sprintf("%s(%s);", callback, result)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	// success!
	return http.StatusOK, result
}

func GetHandler(conf *Config) *martini.ClassicMartini {

	m := martini.Classic()

	m.Get(trimSuffix(conf.Prefix, "/")+"/", NordHandler)

	// shared HTTP client
	client := &http.Client{}

	// shared Redis connection
	conn, err := getRedisConnection(conf)
	if err != nil {
		os.Exit(1)
	}

	serviceMap := &ServiceMap{
		client: client,
		config: conf,
		conn:   conn,
	}
	m.Map(serviceMap)
	return m
}

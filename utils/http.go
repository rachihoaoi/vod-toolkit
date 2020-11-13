package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpRequestConfig struct {
	Client  *http.Client
	Url     string
	Method  string
	Header  map[string]string
	Payload interface{}
	TagName string
}

func DoHttpRequest(config *HttpRequestConfig) (b []byte, respHeader http.Header, err error) {
	var body io.Reader
	var request *http.Request
	switch config.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if config.Payload != nil {
			bytesData, _ := json.Marshal(config.Payload)
			body = bytes.NewBuffer(bytesData)
		}
	case http.MethodGet:
		if config.Payload != nil {
			var tagName string
			var arr = make([]string, 0)
			arr = append(arr, config.Url, "?")
			if config.TagName != "" {
				tagName = config.TagName
			}
			for k, v := range ConvertStructToMap(config.Payload, tagName, true) {
				arr = append(arr, fmt.Sprintf("%s=%v", k, v), "&")
			}
			config.Url = Concatenate(arr)
		}
	}
	if body == nil {
		fmt.Printf("\033[32m[%s]\033[0m \033[32m[%s]\033[0m %s\n", "REQUEST_URL", config.Method, config.Url)
	} else {
		fmt.Printf("\033[32m[%s]\033[0m \033[32m[%s]\033[0m %s\n\033[32m[PAYLOAD]\033[0m: %v \n", "REQUEST_URL", config.Method, config.Url, body)
	}

	if request, err = http.NewRequest(config.Method, config.Url, body); err != nil {
		return
	}
	for k, v := range config.Header {
		request.Header.Set(k, v)
	}
	resp, err := config.Client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	b, err = ioutil.ReadAll(resp.Body)
	return b, resp.Header, err
}

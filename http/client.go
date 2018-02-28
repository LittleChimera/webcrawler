package http

import (
	"io/ioutil"
	"net/http"
)

type Client interface {
	Get(string) string
}

type SimpleClient struct{}

func (s *SimpleClient) Get(url string) string {
	response, err := http.Get(url)
	if err != nil {
		return ""
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}

	defer response.Body.Close()
	return string(body)
}

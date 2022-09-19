package http_helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	neturl "net/url"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func Post(url string, header map[string]string, request interface{}, log ...log.Logger) ([]byte, error) {

	jsonReq, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonReq))

	if err != nil {
		return nil, err
	}

	for k, h := range header {
		req.Header.Add(k, h)
	}

	client := http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)

	for _, v := range log {
		_ = level.Info(v).Log("url", url)
		_ = level.Info(v).Log("request", string(jsonReq))
		_ = level.Info(v).Log("response", string(bodyBytes))
	}

	return bodyBytes, err
}

func Get(url string, header map[string]string, queryString map[string]string, log ...log.Logger) ([]byte, error) {

	q := neturl.Values{}
	for k, v := range queryString {
		q.Add(k, v)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", url, q.Encode()), nil)

	if err != nil {
		return nil, err
	}

	for k, h := range header {
		req.Header.Add(k, h)
	}

	client := http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)

	for _, v := range log {
		_ = level.Info(v).Log("url", url)
		_ = level.Info(v).Log("request", q.Encode())
		_ = level.Info(v).Log("response", string(bodyBytes))
	}

	return bodyBytes, err
}

func Patch(url string, header map[string]string, request interface{}, log ...log.Logger) ([]byte, error) {

	jsonReq, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonReq))

	if err != nil {
		return nil, err
	}

	for k, h := range header {
		req.Header.Add(k, h)
	}

	client := http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)

	for _, v := range log {
		_ = level.Info(v).Log("url", url)
		_ = level.Info(v).Log("request", string(jsonReq))
		_ = level.Info(v).Log("response", string(bodyBytes))
	}

	return bodyBytes, err
}

// func getString(url string, header map[string]string, queryString string) ([]byte, error) {

// 	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", url, queryString), nil)

// 	if err != nil {
// 		return nil, err
// 	}

// 	for k, h := range header {
// 		req.Header.Add(k, h)
// 	}

// 	client := http.Client{}
// 	response, err := client.Do(req)

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer response.Body.Close()
// 	bodyBytes, err := ioutil.ReadAll(response.Body)

// 	return bodyBytes, err
// }

// func update(url string, header map[string]string, request interface{}) ([]byte, error) {
// 	jsonReq, err := json.Marshal(request)

// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonReq))

// 	if err != nil {
// 		return nil, err
// 	}

// 	for k, h := range header {
// 		req.Header.Add(k, h)
// 	}

// 	client := http.Client{}
// 	response, err := client.Do(req)

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer response.Body.Close()
// 	bodyBytes, err := ioutil.ReadAll(response.Body)

// 	return bodyBytes, err
// }

// func delete(url string, header map[string]string, request interface{}) ([]byte, error) {
// 	jsonReq, err := json.Marshal(request)

// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jsonReq))

// 	if err != nil {
// 		return nil, err
// 	}

// 	for k, h := range header {
// 		req.Header.Add(k, h)
// 	}

// 	client := http.Client{}
// 	response, err := client.Do(req)

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer response.Body.Close()
// 	bodyBytes, err := ioutil.ReadAll(response.Body)

// 	return bodyBytes, err
// }

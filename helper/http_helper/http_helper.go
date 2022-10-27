package http_helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	neturl "net/url"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func Post(url string, header map[string]string, request interface{}, log log.Logger) ([]byte, error) {
	var (
		//curl      string
		jsonReq   []byte
		bodyBytes []byte
		err       error
	)

	defer func() {
		_ = level.Info(log).Log("url", url)
		_ = level.Info(log).Log("request", string(jsonReq))
		_ = level.Info(log).Log("response", string(bodyBytes))
		//_ = level.Info(log).Log("curl", curl)
	}()

	jsonReq, err = json.Marshal(request)
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

	//curl, _ = GetCurl(req)
	client := http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	bodyBytes, err = ioutil.ReadAll(response.Body)
	return bodyBytes, err
}

func Get(url string, header map[string]string, queryString map[string]string, log log.Logger) ([]byte, error) {
	var (
		jsonReq   []byte
		bodyBytes []byte
		err       error
	)

	defer func() {
		_ = level.Info(log).Log("url", url)
		_ = level.Info(log).Log("request", string(jsonReq))
		_ = level.Info(log).Log("response", string(bodyBytes))
	}()

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
	bodyBytes, err = ioutil.ReadAll(response.Body)

	return bodyBytes, err
}

func Patch(url string, header map[string]string, request interface{}, log log.Logger) ([]byte, error) {
	var (
		jsonReq   []byte
		bodyBytes []byte
		err       error
	)

	defer func() {
		_ = level.Info(log).Log("url", url)
		_ = level.Info(log).Log("request", string(jsonReq))
		_ = level.Info(log).Log("response", string(bodyBytes))
	}()

	jsonReq, err = json.Marshal(request)

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
	bodyBytes, err = ioutil.ReadAll(response.Body)

	return bodyBytes, err
}

func Delete(url string, header map[string]string, request interface{}, log log.Logger) ([]byte, error) {
	var (
		jsonReq   []byte
		bodyBytes []byte
		err       error
	)

	defer func() {
		_ = level.Info(log).Log("url", url)
		_ = level.Info(log).Log("request", string(jsonReq))
		_ = level.Info(log).Log("response", string(bodyBytes))
	}()

	jsonReq, err = json.Marshal(request)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jsonReq))

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
	bodyBytes, err = ioutil.ReadAll(response.Body)

	return bodyBytes, err
}

func GetCurl(req *http.Request) (string, error) {
	curl := []string{"curl -X", escape(req.Method)}
	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return "", err
		}
		req.Body = closer{bytes.NewBuffer(body)}
		bodyEscaped := escape(string(body))
		curl = append(curl, "-d", bodyEscaped)
	}

	for k, v := range req.Header {
		curl = append(curl, "-H", escape(fmt.Sprintf("%s: %s", k, strings.Join(v, " "))))
	}

	curl = append(curl, escape(req.URL.String()))
	return strings.Join(curl, " "), nil
}

func escape(str string) string {
	escape := strings.ReplaceAll(str, `'`, `'\''`)
	return fmt.Sprint("'", escape, "'")
}

type closer struct {
	io.Reader
}

func (closer) Close() error { return nil }

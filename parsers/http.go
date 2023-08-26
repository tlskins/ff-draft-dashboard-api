package parsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

func NewHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout:    60 * time.Second,
			DisableCompression: true,
		},
		Timeout: 60 * time.Second,
	}
}

func HttpRequest(client *http.Client, method, url string, headers map[string][]string, data, out interface{}) error {
	var body io.Reader
	if data != nil {
		if b, ok := data.([]byte); ok {
			body = bytes.NewReader(b)
		} else {
			if b, err := json.Marshal(data); err == nil {
				body = bytes.NewReader(b)
			} else {
				return errors.New(fmt.Sprintf("Error marshaling request data: %v", err.Error()))
			}
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return errors.New(fmt.Sprintf("Error building http request: %v", err.Error()))
	}
	for header, values := range headers {
		req.Header[header] = values
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("Error sending http request: %v", err.Error()))
	}

	if out != nil {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("Error reading response body: %v", err.Error()))
		}
		if len(b) > 0 {
			// fmt.Println(string(b))
			return json.Unmarshal(b, out)
		}
	}

	return nil
}

func HttpHtmlRequest(client *http.Client, method, url string, headers map[string][]string, data interface{}) (out string, err error) {
	var body io.Reader
	if data != nil {
		if b, ok := data.([]byte); ok {
			body = bytes.NewReader(b)
		} else {
			if b, jErr := json.Marshal(data); jErr == nil {
				body = bytes.NewReader(b)
			} else {
				err = errors.New(fmt.Sprintf("Error marshaling request data: %v", jErr.Error()))
				return
			}
		}
	}

	req, hErr := http.NewRequest(method, url, body)
	if hErr != nil {
		err = errors.New(fmt.Sprintf("Error building http request: %v", hErr.Error()))
		return
	}
	for header, values := range headers {
		req.Header[header] = values
	}

	resp, cErr := client.Do(req)
	if cErr != nil {
		err = errors.New(fmt.Sprintf("Error sending http request: %v", cErr.Error()))
		return
	}

	b, iErr := ioutil.ReadAll(resp.Body)
	if iErr != nil {
		err = errors.New(fmt.Sprintf("Error reading response body: %v", iErr.Error()))
		return
	}
	out = string(b)
	return
}
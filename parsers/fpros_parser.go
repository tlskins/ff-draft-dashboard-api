package parsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

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

const FProsApiUrl = "https://www.fantasypros.com/nfl/rankings/ppr.php"

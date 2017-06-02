package synthetics

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPClient is the interface to the HTTP clients that a Client can
// use.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpClientWithRetries struct {
	client  HTTPClient
	retries uint
}

func newHTTPClientWithRetries(client HTTPClient, retries uint) *httpClientWithRetries {
	return &httpClientWithRetries{client: client, retries: retries}
}

func (h *httpClientWithRetries) Do(req *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if req.Body != nil {
		body := &bytes.Buffer{}
		if _, err := io.Copy(body, req.Body); err != nil {
			return nil, err
		}
		bodyBytes = body.Bytes()
	}

	var response *http.Response
	for i := uint(0); i < h.retries; i++ {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		var err error
		response, err = h.client.Do(req)
		if err != nil {
			return nil, err
		}

		if response.StatusCode == http.StatusTooManyRequests {
			if i != h.retries-1 {
				time.Sleep((1 << i) * time.Second)
			}
			continue
		}

		return response, nil
	}

	return response, nil
}

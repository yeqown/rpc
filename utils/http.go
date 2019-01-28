package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// RequestHTTP ...
func RequestHTTP(serverAddr string, data []byte) ([]byte, error) {
	dataEncoded := url.QueryEscape(fmt.Sprintf("%s", data))
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://%s?data=%s", serverAddr, dataEncoded),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest got err: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.DefaultClient.Do(req) got err: %v", err)
	}

	// read from response body
	byts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("http.DefaultClient.Do(req) got err: %v", err)
	}
	defer resp.Body.Close()
	return byts, nil
}

// ResponseHTTP ... for server side to write response to the client
func ResponseHTTP(w http.ResponseWriter, data []byte, logResponse bool) {
	if _, err := io.WriteString(w, string(data)); err != nil {
		panic(err)
	}

	if logResponse {
		log.Printf("[HTTP] response to client: %s", data)
	}
}

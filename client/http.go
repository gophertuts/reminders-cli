package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// reminderBody represents reminder request body
type reminderBody struct {
	ID       string        `json:"id"`
	Title    string        `json:"title"`
	Message  string        `json:"message"`
	Duration time.Duration `json:"duration"`
}

// httpRoundTripper represents the HTTP interceptor for the CLI HTTP client
type httpRoundTripper struct {
	healthURI         string
	originalTransport http.RoundTripper
}

func (roundTripper *httpRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := http.Get(roundTripper.healthURI)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend api is not available")
	}
	res, err = roundTripper.originalTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// HTTPClient represents the HTTP client which communicates with reminders backend API
type HTTPClient struct {
	client     *http.Client
	BackendURI string
}

// NewHTTPClient creates a new instance of HTTPClient
func NewHTTPClient(uri string) HTTPClient {
	roundTripper := &httpRoundTripper{
		healthURI:         uri + "/health",
		originalTransport: http.DefaultTransport,
	}
	return HTTPClient{
		BackendURI: uri,
		client: &http.Client{
			Transport: roundTripper,
		},
	}
}

// Create calls the create API endpoint
func (c HTTPClient) Create(title, message string, duration time.Duration) ([]byte, error) {
	requestBody := reminderBody{
		Title:    title,
		Message:  message,
		Duration: duration,
	}
	return c.apiCall(
		http.MethodPost,
		"/reminders",
		&requestBody,
		http.StatusCreated,
	)
}

// Edit calls the edit API endpoint
func (c HTTPClient) Edit(id string, title, message string, duration time.Duration) ([]byte, error) {
	requestBody := reminderBody{
		ID:       id,
		Title:    title,
		Message:  message,
		Duration: duration,
	}
	return c.apiCall(
		http.MethodPatch,
		"/reminders"+id,
		&requestBody,
		http.StatusOK,
	)
}

// Fetch calls the fetch API endpoint
func (c HTTPClient) Fetch(ids []string) ([]byte, error) {
	idsSet := strings.Join(ids, ",")
	return c.apiCall(
		http.MethodGet,
		"/reminders/"+idsSet,
		nil,
		http.StatusOK,
	)
}

// Delete calls the delete API endpoint
func (c HTTPClient) Delete(ids []string) error {
	idsSet := strings.Join(ids, ",")
	_, err := c.apiCall(
		http.MethodDelete,
		"/reminders/"+idsSet,
		nil,
		http.StatusNoContent,
	)
	return err
}

// apiCall makes a new backend api call
func (c HTTPClient) apiCall(method, path string, body interface{}, resCode int) ([]byte, error) {
	bs, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, c.BackendURI+path, bytes.NewReader(bs))
	if err != nil {
		return []byte{}, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode != resCode {
		return []byte{}, fmt.Errorf(
			"expected response code: %d, got: %d",
			resCode,
			res.StatusCode,
		)
	}

	return c.readBody(res.Body)
}

func (c HTTPClient) readBody(b io.Reader) ([]byte, error) {
	bs, err := ioutil.ReadAll(b)
	if err != nil {
		return []byte{}, err
	}

	var buff bytes.Buffer
	err = json.Indent(&buff, bs, "", "\t")
	if err != nil {
		return []byte{}, err
	}

	return buff.Bytes(), nil
}

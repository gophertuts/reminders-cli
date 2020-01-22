package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gophertuts/reminders-cli/server/models"
)

// HTTPClient represents the HTTP client for communicating with the notifier server
type HTTPClient struct {
	notifierURI string
	client      *http.Client
}

// NewHTTPClient creates a new HTTP client instance
func NewHTTPClient(uri string) HTTPClient {
	return HTTPClient{
		notifierURI: uri,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Notify pushes a given reminder to the notifier service
// if the reminder is nil, means the record must be retried
func (c HTTPClient) Notify(reminder models.Reminder) (time.Duration, error) {
	var notifierResponse struct {
		ActivationType  string `json:"activationType"`
		ActivationValue string `json:"activationValue"`
	}
	bs, err := json.Marshal(reminder)
	if err != nil {
		e := models.WrapError("could not marshal json", err)
		return 0, e
	}

	res, err := c.client.Post(
		c.notifierURI+"/notify",
		"application/json",
		bytes.NewReader(bs),
	)
	if err != nil {
		e := models.WrapError("notifier service is not available", err)
		return 0, e
	}

	err = json.NewDecoder(res.Body).Decode(&notifierResponse)
	if err != nil && err != io.EOF {
		e := models.WrapError("could not decode notifier response", err)
		return 0, e
	}
	if notifierResponse.ActivationType != "replied" {
		return 0, nil
	}

	d, err := time.ParseDuration(notifierResponse.ActivationValue)
	if err != nil && d != 0 {
		e := models.WrapError("could not parse notifier duration", err)
		return 0, e
	}
	return d, nil
}

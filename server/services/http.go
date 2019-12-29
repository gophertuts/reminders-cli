package services

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
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
// if nil is returned, means the record must be retried
func (c HTTPClient) Notify(reminder models.Reminder) (*models.Reminder, time.Duration) {
	res, err := c.client.Get(c.notifierURI + "/health")
	if err != nil {
		log.Printf("notifier service is not available: %v\n", err)
		return &reminder, 0
	}

	var notifierResponse struct {
		ActivationType  string `json:"activationType"`
		ActivationValue string `json:"activationValue"`
	}
	bs, err := json.Marshal(reminder)
	if err != nil {
		log.Printf("could not marshal json: %v\n", err)
		// means FORGET ABOUT THIS RECORD
		return nil, 0
	}
	res, err = c.client.Post(
		c.notifierURI+"/notify",
		"application/json",
		bytes.NewReader(bs),
	)
	if err != nil {
		log.Printf("something went wrong with the notifier: %v\n", err)
		return &reminder, 0
	}
	err = json.NewDecoder(res.Body).Decode(&notifierResponse)
	if err != nil && err != io.EOF {
		log.Printf("could not decode notifier response: %v\n", err)
	}
	if notifierResponse.ActivationType == "replied" {
		d, err := time.ParseDuration(notifierResponse.ActivationValue)
		if err != nil {
			log.Printf("could not parse input duration: %v\n", err)
		}
		return &reminder, d
	}
	return nil, 0
}

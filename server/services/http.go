package services

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
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
			Timeout: 20 * time.Second,
		},
	}
}

// NotificationResponse represents OS notification response for background notifier
type NotificationResponse struct {
	completed bool
	duration  time.Duration
}

// Notify pushes a given reminder to the notifier service
// if the reminder is nil, means the record must be retried
func (c HTTPClient) Notify(reminder models.Reminder) (NotificationResponse, error) {
	var notifierResponse struct {
		ActivationType  string `json:"activationType"`
		ActivationValue string `json:"activationValue"`
	}
	bs, err := json.Marshal(reminder)
	if err != nil {
		e := models.WrapError("could not marshal json", err)
		return NotificationResponse{}, e
	}

	res, err := c.client.Post(
		c.notifierURI+"/notify",
		"application/json",
		bytes.NewReader(bs),
	)
	if err != nil {
		e := models.WrapError("notifier service is not available", err)
		return NotificationResponse{}, e
	}
	err = json.NewDecoder(res.Body).Decode(&notifierResponse)
	if err != nil && err != io.EOF {
		e := models.WrapError("could not decode notifier response", err)
		return NotificationResponse{}, e
	}

	t := notifierResponse.ActivationType
	v := notifierResponse.ActivationValue
	if t == "closed" {
		return NotificationResponse{completed: true}, nil
	}

	d, err := time.ParseDuration(v)
	if err != nil && d != 0 {
		e := models.WrapError("could not parse notifier duration", err)
		return NotificationResponse{}, e
	}
	if d == 0 {
		return NotificationResponse{}, errors.New("notification duration must be > 0s")
	}
	return NotificationResponse{duration: d}, nil
}

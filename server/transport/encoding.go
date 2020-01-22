package transport

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gophertuts/reminders-cli/server/models"
)

// SendJSON sends a json response to the client
func SendJSON(w http.ResponseWriter, response interface{}, code int) {
	encoder := jsonEncoder(w, code)
	if err := encoder.Encode(response); err != nil {
		log.Printf("could not encode error: %v", err)
	}
}

// SendError sends a json error to the client
func SendError(w http.ResponseWriter, err error, code int) {
	encoder := jsonEncoder(w, code)
	e := toHTTPError(err)
	if err := encoder.Encode(e); err != nil {
		log.Printf("could not encode error: %v", err)
	}
}

// jsonEncoder creates a new json encoder
func jsonEncoder(w http.ResponseWriter, code int) *json.Encoder {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w)
}

// toHTTPError converts an error to HTTPError
func toHTTPError(err error) models.HTTPError {
	return models.HTTPError{Message: err.Error()}
}

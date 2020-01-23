package transport

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gophertuts/reminders-cli/server/models"
)

const (
	notFoundErrType         = "resource_not_found_error"
	dataValidationErrType   = "data_validation_error"
	formatValidationErrType = "format_validation_error"
	invalidJSONErrType      = "invalid_json_error"
	serviceErrType          = "service_error"
)

// SendJSON sends a json response to the client
func SendJSON(w http.ResponseWriter, response interface{}, code int) {
	encoder := jsonEncoder(w, code)
	if err := encoder.Encode(response); err != nil {
		log.Printf("could not encode error: %v", err)
	}
}

// SendError sends a json error to the client
func SendError(w http.ResponseWriter, err error) {
	e := toHTTPError(err)
	encoder := jsonEncoder(w, e.Code)
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
	resErr := models.HTTPError{Message: err.Error()}
	switch e := err.(type) {
	case models.HTTPError:
		return e
	case models.NotFoundError:
		resErr.Code = http.StatusNotFound
		resErr.Type = notFoundErrType
	case models.FormatValidationError:
		resErr.Code = http.StatusBadRequest
		resErr.Type = formatValidationErrType
	case models.DataValidationError:
		resErr.Code = http.StatusBadRequest
		resErr.Type = dataValidationErrType
	case models.InvalidJSONError:
		resErr.Code = http.StatusBadRequest
		resErr.Type = invalidJSONErrType
	default:
		resErr.Code = http.StatusInternalServerError
		resErr.Type = serviceErrType
		resErr.Message = "Internal Server Error"
	}
	log.Printf("error: %v", err)
	return resErr
}

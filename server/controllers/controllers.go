package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// checkHTTPMethod checks whether the given HTTP method matches the one from request
func checkHTTPMethod(w http.ResponseWriter, r *http.Request, method string) {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// jsonEncode encodes data into json
func jsonEncode(w io.Writer, v interface{}) {
	if err := json.NewEncoder(w).Encode(&v); err != nil {
		log.Fatalf("could not encode json: %v", err)
	}
}

// jsonDecode decodes json into data
func jsonDecode(r io.Reader, v interface{}) {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		log.Fatalf("could not decode json: %v", err)
	}
}

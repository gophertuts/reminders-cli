package controllers

import "net/http"

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

package controllers

import (
	"net/http"
)

type deleter interface {
	Delete(ids []int)
}

func deleteReminders(service deleter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(w, r, http.MethodDelete)
		var body struct {
			IDs []int `json:"ids"`
		}
		jsonDecode(r.Body, &body)
		service.Delete(body.IDs)
		w.WriteHeader(http.StatusNoContent)
	})
}

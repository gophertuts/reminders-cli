package controllers

import (
	"github.com/gophertuts/reminders-cli/server/transport"
	"net/http"
)

type deleter interface {
	Delete(ids []int) error
}

func deleteReminders(service deleter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids := parseIDsParam(r.Context())
		err := service.Delete(ids)
		if err != nil {
			transport.SendError(w, err, http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

package controllers

import (
	"net/http"

	"github.com/gophertuts/reminders-cli/server/models"
)

type fetcher interface {
	Fetch(ids []int) []models.Reminder
}

func fetchReminders(service fetcher) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids := parseIDsParam(r.Context())
		reminders := service.Fetch(ids)
		jsonEncode(w, reminders, http.StatusOK)
	})
}

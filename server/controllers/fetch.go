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
		checkHTTPMethod(w, r, http.MethodPost)
		var body struct {
			IDs []int `json:"ids"`
		}
		jsonDecode(r.Body, &body)
		reminders := service.Fetch(body.IDs)
		jsonEncode(w, reminders)
	})
}

package controllers

import (
	"net/http"
	"time"

	"github.com/gophertuts/reminders-cli/server/models"
	"github.com/gophertuts/reminders-cli/server/services"
)

type creator interface {
	Create(reminderBody services.ReminderCreateBody) models.Reminder
}

func createReminder(service creator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(w, r, http.MethodPost)
		var body struct {
			Title    string        `json:"title"`
			Message  string        `json:"message"`
			Duration time.Duration `json:"duration"`
		}
		jsonDecode(r.Body, &body)
		reminder := service.Create(services.ReminderCreateBody{
			Title:    body.Title,
			Message:  body.Message,
			Duration: body.Duration,
		})
		w.WriteHeader(http.StatusCreated)
		jsonEncode(w, reminder)
	})
}

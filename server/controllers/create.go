package controllers

import (
	"net/http"
	"time"

	"github.com/gophertuts/reminders-cli/server/models"
)

// ReminderCreateBody represents HTTP body for creating a reminder
type ReminderCreateBody struct {
	Title    string        `json:"title"`
	Message  string        `json:"message"`
	Duration time.Duration `json:"duration"`
}

type creator interface {
	Create(reminderBody ReminderCreateBody) models.Reminder
}

func createReminder(service creator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(w, r, http.MethodPost)
		var body ReminderCreateBody
		jsonDecode(r.Body, &body)
		reminder := service.Create(body)
		w.WriteHeader(http.StatusCreated)
		jsonEncode(w, reminder)
	})
}

package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gophertuts/reminders-cli/server/models"
)

// ReminderEditBody represents HTTP body editing a reminder
type ReminderEditBody struct {
	ID       int           `json:"id"`
	Title    string        `json:"title"`
	Message  string        `json:"message"`
	Duration time.Duration `json:"duration"`
}

type editor interface {
	Edit(reminderBody ReminderEditBody) (models.Reminder, error)
}

func editReminder(service editor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(w, r, http.MethodPut)
		var body ReminderEditBody
		jsonDecode(r.Body, &body)
		reminder, err := service.Edit(body)
		if err != nil {
			log.Fatalf("could not edit reminder: %v", err)
		}
		jsonEncode(w, reminder)
	})
}

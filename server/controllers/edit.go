package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gophertuts/reminders-cli/server/models"
	"github.com/gophertuts/reminders-cli/server/services"
)

type editor interface {
	Edit(reminderBody services.ReminderEditBody) (models.Reminder, error)
}

func editReminder(service editor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(w, r, http.MethodPut)
		var body struct {
			ID       int           `json:"id"`
			Title    string        `json:"title"`
			Message  string        `json:"message"`
			Duration time.Duration `json:"duration"`
		}
		jsonDecode(r.Body, &body)
		reminder, err := service.Edit(services.ReminderEditBody{
			ID:       body.ID,
			Title:    body.Title,
			Message:  body.Message,
			Duration: body.Duration,
		})
		if err != nil {
			log.Fatalf("could not edit reminder: %v", err)
		}
		jsonEncode(w, reminder)
	})
}

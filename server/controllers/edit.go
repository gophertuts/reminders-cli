package controllers

import (
	"encoding/json"
	"github.com/gophertuts/reminders-cli/server/transport"
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
		id, err := parseIDParam(r.Context())
		if err != nil {
			transport.SendError(w, err)
			return
		}
		var body struct {
			Title    string        `json:"title"`
			Message  string        `json:"message"`
			Duration time.Duration `json:"duration"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			transport.SendError(w, models.InvalidJSONError{Message: err.Error()})
			return
		}
		reminder, err := service.Edit(services.ReminderEditBody{
			ID:       id,
			Title:    body.Title,
			Message:  body.Message,
			Duration: body.Duration,
		})
		if err != nil {
			transport.SendError(w, err)
			return
		}
		transport.SendJSON(w, reminder, http.StatusOK)
	})
}

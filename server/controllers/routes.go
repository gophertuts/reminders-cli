package controllers

import (
	"net/http"

	"github.com/gophertuts/reminders-cli/server/middleware"
)

// HTTP params
const (
	idParamName  = "id"
	idsParamName = "ids"
	idParam      = `{` + idParamName + `}:^[0-9]+$`
	idsParam     = `{` + idsParamName + `}:[0-9]+(,[0-9]+)*`
)

// RemindersService represents the Reminders service
type RemindersService interface {
	creator
	editor
	fetcher
	deleter
}

// RouterConfig represents router specific configuration
type RouterConfig struct {
	Service RemindersService
}

// NewRouter creates a new server (backend) application router
func NewRouter(cfg RouterConfig) http.Handler {
	r := RegexpMux{}
	m := middleware.New(
		middleware.HTTPLogger,
	)
	r.Get("/health", m.Then(health()))
	r.Post("/reminders", m.Then(createReminder(cfg.Service)))
	r.Get("/reminders/"+idsParam, m.Then(fetchReminders(cfg.Service)))
	r.Delete("/reminders/"+idsParam, m.Then(deleteReminders(cfg.Service)))
	r.Patch("/reminders/"+idParam, m.Then(editReminder(cfg.Service)))
	return r
}

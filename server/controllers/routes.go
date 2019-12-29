package controllers

import (
	"net/http"

	"github.com/gophertuts/reminders-cli/server/middleware"
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
	r := http.NewServeMux()
	m := middleware.New(
		middleware.HTTPLogger,
	)
	r.HandleFunc("/health", health)
	r.Handle("/reminders/create", m.Then(createReminder(cfg.Service)))
	r.Handle("/reminders/edit", m.Then(editReminder(cfg.Service)))
	r.Handle("/reminders/fetch", m.Then(fetchReminders(cfg.Service)))
	r.Handle("/reminders/delete", m.Then(deleteReminders(cfg.Service)))
	return r
}

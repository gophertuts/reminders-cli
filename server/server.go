package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gophertuts/reminders-cli/server/controllers"
	"github.com/gophertuts/reminders-cli/server/repositories"
	"github.com/gophertuts/reminders-cli/server/services"
)

// AppConfig represents application configuration
type AppConfig struct {
	Addr        string
	NotifierURI string
	DB          *repositories.DB
}

// App represents the server (backend) application
type App struct {
	cfg                AppConfig
	Server             *http.Server
	DB                 *repositories.DB
	BackgroundSaver    *services.BackgroundSaver
	BackgroundNotifier *services.BackgroundNotifier
}

// New initializes and creates a new server App
func New(c AppConfig) *App {
	repo := repositories.NewReminder(c.DB)
	service := services.NewReminders(repo)
	saver := services.NewBackgroundSaver(service)
	notifier := services.NewBackgroundNotifier(c.NotifierURI, service)
	routerCfg := controllers.RouterConfig{
		Service: service,
	}
	router := controllers.NewRouter(routerCfg)
	return &App{
		cfg: c,
		DB:  c.DB,
		Server: &http.Server{
			Addr:    c.Addr,
			Handler: router,
		},
		BackgroundSaver:    saver,
		BackgroundNotifier: notifier,
	}
}

// Start starts the initialized server (backend) application
func (a *App) Start() error {
	log.Printf("application started on address %s\n", a.cfg.Addr)
	go a.BackgroundSaver.Start()
	go a.BackgroundNotifier.Start()
	err := a.Server.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Println("http server is closed")
		return nil
	}
	return err
}

// Kill gracefully kills the server (backend) application
func (a *App) Kill() {
	timeout := 2 * time.Second
	ch := make(chan struct{})
	a.BackgroundSaver.Stop()
	a.BackgroundNotifier.Stop()
	go func() {
		log.Println("shutting down the http server")
		if err := a.Server.Shutdown(context.Background()); err != nil {
			log.Fatalf("error on server shutdown: %v", err)
		}
		close(ch)
	}()
	select {
	case <-ch:
		log.Println("application was shut down")
	case <-time.After(timeout):
		log.Printf("could not shut down in: %v\n", timeout)
	}
}

// Killer represents applications (servers) which can be closed (killed) gracefully
type Killer interface {
	Kill()
}

// ListenForSignals listens for OS signals and responds to the provided ones
func ListenForSignals(signals []os.Signal, servers ...Killer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)
	sig := <-c
	log.Printf("received shutdown signal: %v\n", sig.String())

	for _, server := range servers {
		server.Kill()
	}
	os.Exit(0)
}

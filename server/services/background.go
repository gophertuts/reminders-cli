package services

import (
	"log"
	"time"

	"github.com/gophertuts/reminders-cli/server/models"
)

type saver interface {
	save()
}

// BackgroundSaver represents the reminder background saver
type BackgroundSaver struct {
	ticker  *time.Ticker
	service saver
}

// NewBackgroundSaver creates a new instance of BackgroundSaver
func NewBackgroundSaver(service saver) *BackgroundSaver {
	ticker := time.NewTicker(30 * time.Second)
	return &BackgroundSaver{
		ticker:  ticker,
		service: service,
	}
}

// Start starts the created Watcher
func (s *BackgroundSaver) Start() {
	log.Println("background saver started")
	for {
		select {
		case <-s.ticker.C:
			s.service.save()
		}
	}
}

// Stop stops the created Watcher
func (s *BackgroundSaver) Stop() {
	s.ticker.Stop()
	s.service.save()
	log.Println("background saver stopped")
}

// HTTPNotifierClient represents the HTTP client for communicating with the notifier server
type HTTPNotifierClient interface {
	Notify(reminder models.Reminder) (*models.Reminder, time.Duration)
}

type snapshotManager interface {
	snapshot() Snapshot
	snapshotGrooming(notifiedReminders ...models.Reminder)
	retry(reminder models.Reminder, duration time.Duration)
}

// BackgroundNotifier represents the reminder background saver
type BackgroundNotifier struct {
	ticker  *time.Ticker
	service snapshotManager
	Client  HTTPNotifierClient
}

// NewBackgroundNotifier creates a new instance of BackgroundNotifier
func NewBackgroundNotifier(notifierURI string, service snapshotManager) *BackgroundNotifier {
	ticker := time.NewTicker(1 * time.Second)
	httpClient := NewHTTPClient(notifierURI)
	return &BackgroundNotifier{
		ticker:  ticker,
		service: service,
		Client:  httpClient,
	}
}

// Start starts the created Watcher
func (s *BackgroundNotifier) Start() {
	log.Println("background notifier started")
	for {
		select {
		case <-s.ticker.C:
			snapshot := s.service.snapshot()
			notified := map[int]models.Reminder{}
			for _, reminder := range snapshot.UnCompleted {
				reminderTick := reminder.ModifiedAt.Add(reminder.Duration).UnixNano()
				nowTick := time.Now().UnixNano()
				deltaTick := time.Now().Add(time.Second).UnixNano()
				if reminderTick > nowTick && reminderTick < deltaTick {
					go func(r models.Reminder) {
						retry, d := s.Client.Notify(r)
						if retry != nil {
							s.service.retry(r, d)
						} else {
							s.service.snapshotGrooming(r)
						}
					}(reminder)
					notified[reminder.ID] = reminder
				}
			}
			if len(notified) > 0 {
				log.Printf("notified: %d record(s)\n", len(notified))
			}
		}
	}
}

// Stop stops the created Watcher
func (s *BackgroundNotifier) Stop() {
	s.ticker.Stop()
	log.Println("background notifier stopped")
}

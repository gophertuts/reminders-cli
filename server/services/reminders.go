package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gophertuts/reminders-cli/server/controllers"
	"github.com/gophertuts/reminders-cli/server/models"
)

const (
	retryPeriod = time.Minute
)

// ReminderRepository represents the Reminder repository
type ReminderRepository interface {
	Save([]models.Reminder)
	Filter(filterFn func(reminder models.Reminder) bool) (map[int]models.Reminder, map[int]int)
	NextID() int
}

// Snapshot represents current service in memory state
type Snapshot struct {
	All           map[int]models.Reminder
	UnCompleted   map[int]models.Reminder
	OriginalOrder map[int]int
}

// Reminders represents the Reminders service
type Reminders struct {
	repo     ReminderRepository
	Snapshot Snapshot
}

// NewReminders creates a new instance of Reminders service
func NewReminders(repo ReminderRepository) Reminders {
	all, originalOrder := repo.Filter(nil)
	unCompleted, _ := repo.Filter(func(r models.Reminder) bool {
		return r.ModifiedAt.Add(r.Duration).UnixNano() > time.Now().UnixNano()
	})
	return Reminders{
		repo: repo,
		Snapshot: Snapshot{
			All:           all,
			UnCompleted:   unCompleted,
			OriginalOrder: originalOrder,
		},
	}
}

// Create creates a new Reminder
func (r Reminders) Create(reminderBody controllers.ReminderCreateBody) models.Reminder {
	reminder := models.Reminder{
		ID:         r.repo.NextID(),
		Title:      reminderBody.Title,
		Message:    reminderBody.Message,
		Duration:   reminderBody.Duration,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
	r.Snapshot.All[reminder.ID] = reminder
	r.Snapshot.UnCompleted[reminder.ID] = reminder
	r.Snapshot.OriginalOrder[reminder.ID] = len(r.Snapshot.OriginalOrder)
	return reminder
}

// Edit edits a given Reminder
func (r Reminders) Edit(reminderBody controllers.ReminderEditBody) (models.Reminder, error) {
	fmt.Println(reminderBody)
	reminder, ok := r.Snapshot.All[reminderBody.ID]
	if !ok {
		err := fmt.Errorf("could not find reminder with id: %d", reminderBody.ID)
		return models.Reminder{}, err
	}
	if strings.TrimSpace(reminderBody.Title) != "" {
		reminder.Title = reminderBody.Title
	}
	if strings.TrimSpace(reminderBody.Message) != "" {
		reminder.Message = reminderBody.Message
	}
	if reminderBody.Duration > 0 {
		fmt.Println("set duration")
		reminder.Duration = reminderBody.Duration
	}
	reminder.ModifiedAt = time.Now()
	r.Snapshot.All[reminder.ID] = reminder
	if reminder.ModifiedAt.UnixNano() < time.Now().Add(reminderBody.Duration).UnixNano() {
		r.Snapshot.UnCompleted[reminder.ID] = reminder
	} else {
		delete(r.Snapshot.UnCompleted, reminder.ID)
	}
	return reminder, nil
}

// Fetch fetches a list of reminders
func (r Reminders) Fetch(ids []int) []models.Reminder {
	var reminders []models.Reminder
	for _, id := range ids {
		reminder, ok := r.Snapshot.All[id]
		if ok {
			reminders = append(reminders, reminder)
		}
	}
	return reminders
}

// Delete deletes a list of reminders and persists the changes
func (r Reminders) Delete(ids []int) {
	for _, id := range ids {
		delete(r.Snapshot.All, id)
		delete(r.Snapshot.UnCompleted, id)
		delete(r.Snapshot.OriginalOrder, id)
	}
}

// Save saves the current reminders snapshot
func (r Reminders) save() {
	log.Println("saving current snapshot")
	reminders := make([]models.Reminder, len(r.Snapshot.All))
	for id, i := range r.Snapshot.OriginalOrder {
		reminders[i] = r.Snapshot.All[id]
	}
	r.repo.Save(reminders)
}

// GetSnapshot fetches the current service snapshot
func (r Reminders) snapshot() Snapshot {
	return r.Snapshot
}

func (r Reminders) snapshotGrooming(notifiedReminders ...models.Reminder) {
	if len(notifiedReminders) > 0 {
		log.Printf("snapshot grooming: %d record(s)", len(notifiedReminders))
	}
	for _, reminder := range notifiedReminders {
		delete(r.Snapshot.UnCompleted, reminder.ID)
		reminder.Duration = -time.Hour
		r.Snapshot.All[reminder.ID] = reminder
	}
}

func (r Reminders) retry(reminder models.Reminder, d time.Duration) {
	log.Printf("retrying record with id: %d ", reminder.ID)
	reminder.ModifiedAt = time.Now()
	if d <= 0 {
		reminder.Duration = retryPeriod
	} else {
		reminder.Duration = d
	}
	r.Snapshot.All[reminder.ID] = reminder
	r.Snapshot.UnCompleted[reminder.ID] = reminder
}

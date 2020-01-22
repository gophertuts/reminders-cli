package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gophertuts/reminders-cli/server/models"
)

const (
	retryPeriod = time.Minute
)

// RemindersMap represents the data structure for in-memory reminders collection
type RemindersMap map[int]map[int]models.Reminder

// flatten assumes the map has only one value and reads it and retrieves the key and value
func (rMap RemindersMap) flatten(id int) (int, models.Reminder) {
	var index int
	var reminder models.Reminder
	for i, r := range rMap[id] {
		index = i
		reminder = r
	}
	return index, reminder
}

// ReminderRepository represents the Reminder repository
type ReminderRepository interface {
	Save([]models.Reminder) (int, error)
	Filter(filterFn func(reminder models.Reminder) bool) (RemindersMap, error)
	NextID() int
}

// Snapshot represents current service in memory state
type Snapshot struct {
	All         RemindersMap
	UnCompleted RemindersMap
}

// Reminders represents the Reminders service
type Reminders struct {
	repo     ReminderRepository
	Snapshot Snapshot
}

// NewReminders creates a new instance of Reminders service
func NewReminders(repo ReminderRepository) *Reminders {
	return &Reminders{
		repo: repo,
		Snapshot: Snapshot{
			All:         RemindersMap{},
			UnCompleted: RemindersMap{},
		},
	}
}

// Populate populates the reminders service internal state with data from db file
func (s *Reminders) Populate() error {
	all, err := s.repo.Filter(nil)
	if err != nil {
		return models.WrapError("could not get all reminders", err)
	}
	unCompleted, err := s.repo.Filter(func(r models.Reminder) bool {
		return r.ModifiedAt.Add(r.Duration).UnixNano() > time.Now().UnixNano()
	})
	if err != nil {
		return models.WrapError("could not get uncompleted reminders", err)
	}
	s.Snapshot.All = all
	s.Snapshot.UnCompleted = unCompleted
	return nil
}

// ReminderCreateBody represents the model for creating a reminder
type ReminderCreateBody struct {
	Title    string
	Message  string
	Duration time.Duration
}

// Create creates a new Reminder
func (s Reminders) Create(reminderBody ReminderCreateBody) models.Reminder {
	nextID := s.repo.NextID()
	reminder := models.Reminder{
		ID:         nextID,
		Title:      reminderBody.Title,
		Message:    reminderBody.Message,
		Duration:   reminderBody.Duration,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
	index := len(s.Snapshot.All)
	s.Snapshot.All[reminder.ID] = map[int]models.Reminder{index: reminder}
	s.Snapshot.UnCompleted[reminder.ID] = map[int]models.Reminder{index: reminder}
	return reminder
}

// ReminderEditBody represents the model for editing a reminder
type ReminderEditBody struct {
	ID       int
	Title    string
	Message  string
	Duration time.Duration
}

// Edit edits a given Reminder
func (s Reminders) Edit(reminderBody ReminderEditBody) (models.Reminder, error) {
	_, ok := s.Snapshot.All[reminderBody.ID]
	if !ok {
		err := fmt.Errorf("could not find reminder with id: %d", reminderBody.ID)
		return models.Reminder{}, err
	}
	index, reminder := s.Snapshot.All.flatten(reminderBody.ID)
	if strings.TrimSpace(reminderBody.Title) != "" {
		reminder.Title = reminderBody.Title
	}
	if strings.TrimSpace(reminderBody.Message) != "" {
		reminder.Message = reminderBody.Message
	}
	if reminderBody.Duration != 0 {
		reminder.Duration = reminderBody.Duration
	}
	reminder.ModifiedAt = time.Now()
	s.Snapshot.All[reminder.ID] = map[int]models.Reminder{index: reminder}
	if reminder.ModifiedAt.UnixNano() < time.Now().Add(reminderBody.Duration).UnixNano() {
		s.Snapshot.UnCompleted[reminder.ID] = map[int]models.Reminder{index: reminder}
	} else {
		delete(s.Snapshot.UnCompleted, reminder.ID)
	}
	return reminder, nil
}

// Fetch fetches a list of reminders
func (s Reminders) Fetch(ids []int) []models.Reminder {
	reminders := []models.Reminder{}
	for _, id := range ids {
		_, ok := s.Snapshot.All[id]
		if !ok {
			// TODO: return an error
		}
		_, reminder := s.Snapshot.All.flatten(id)
		reminders = append(reminders, reminder)
	}
	return reminders
}

// IDsResponse represents response in form of deleted and not found ids
type IDsResponse struct {
	NotFoundIDs []int
	DeletedIDs  []int
}

// Delete deletes a list of reminders and persists the changes
func (s Reminders) Delete(ids []int) IDsResponse {
	var idsRes IDsResponse
	for _, id := range ids {
		_, ok := s.Snapshot.All[id]
		if !ok {
			idsRes.NotFoundIDs = append(idsRes.NotFoundIDs, id)
		} else {
			idsRes.DeletedIDs = append(idsRes.DeletedIDs, id)
			delete(s.Snapshot.All, id)
			delete(s.Snapshot.UnCompleted, id)
		}
	}
	return idsRes
}

// save saves the current reminders snapshot
func (s Reminders) save() error {
	reminders := make([]models.Reminder, len(s.Snapshot.All))
	for _, reminderMap := range s.Snapshot.All {
		for i, reminder := range reminderMap {
			reminders[i] = reminder
		}
	}

	n, err := s.repo.Save(reminders)
	if err != nil {
		return models.WrapError("could not save snapshot", err)
	}
	if n > 0 {
		log.Printf("successfully saved snapshot: %d reminders", len(reminders))
	}
	return nil
}

// GetSnapshot fetches the current service snapshot
func (s Reminders) snapshot() Snapshot {
	return s.Snapshot
}

// snapshotGrooming clears the current snapshot from notified reminders
func (s Reminders) snapshotGrooming(notifiedReminders ...models.Reminder) {
	if len(notifiedReminders) > 0 {
		log.Printf("snapshot grooming: %d record(s)", len(notifiedReminders))
	}
	for _, reminder := range notifiedReminders {
		delete(s.Snapshot.UnCompleted, reminder.ID)
		reminder.Duration = -time.Hour
		index, _ := s.Snapshot.All.flatten(reminder.ID)
		s.Snapshot.All[reminder.ID] = map[int]models.Reminder{index: reminder}
	}
}

// retry retries a reminder by resetting its duration
func (s Reminders) retry(reminder models.Reminder, d time.Duration) {
	reminder.ModifiedAt = time.Now()
	if d <= 0 {
		reminder.Duration = retryPeriod
	} else {
		reminder.Duration = d
	}
	log.Printf(
		"retrying record with id: %d after %v",
		reminder.ID,
		reminder.Duration.String(),
	)
	index, _ := s.Snapshot.All.flatten(reminder.ID)
	s.Snapshot.All[reminder.ID] = map[int]models.Reminder{index: reminder}
	s.Snapshot.UnCompleted[reminder.ID] = map[int]models.Reminder{index: reminder}
}

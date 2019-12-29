package repositories

import (
	"github.com/gophertuts/reminders-cli/server/models"
)

// Reminder represents the Reminder repository (database layer)
type Reminder struct {
	DB *DB
}

// NewReminder creates a new instance of Reminder repository
func NewReminder(db *DB) Reminder {
	return Reminder{
		DB: db,
	}
}

// Save saves the current snapshot of reminders in the DB
func (r Reminder) Save(reminders []models.Reminder) {
	r.DB.Write(reminders)
}

// Filter filters reminders by a filtering function
func (r Reminder) Filter(filterFn func(reminder models.Reminder) bool) (map[int]models.Reminder, map[int]int) {
	remindersMap := map[int]models.Reminder{}
	originalOrder := map[int]int{}
	for i, reminder := range r.DB.ReadAll() {
		if filterFn == nil || filterFn(reminder) {
			remindersMap[reminder.ID] = reminder
			originalOrder[reminder.ID] = i
		}
	}
	return remindersMap, originalOrder
}

// NextID fetches the next DB AUTOINCREMENT id
func (r Reminder) NextID() int {
	return r.DB.genID()
}

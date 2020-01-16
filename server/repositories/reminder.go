package repositories

import (
	"encoding/json"
	"io"
	"log"

	"github.com/gophertuts/reminders-cli/server/models"
)

// FileDB represents the file database
type FileDB interface {
	io.ReadWriter
	SizeOf() int
	GenerateID() int
}

// Reminder represents the Reminder repository (database layer)
type Reminder struct {
	DB FileDB
}

// NewReminder creates a new instance of Reminder repository
func NewReminder(db FileDB) Reminder {
	return Reminder{
		DB: db,
	}
}

// Save saves the current snapshot of reminders in the DB
func (r Reminder) Save(reminders []models.Reminder) (int, error) {
	bs, err := json.Marshal(reminders)
	if err != nil {
		return 0, err
	}
	n, err := r.DB.Write(bs)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// Filter filters reminders by a filtering function
func (r Reminder) Filter(filterFn func(reminder models.Reminder) bool) (map[int]models.Reminder, map[int]int) {
	bs := make([]byte, r.DB.SizeOf())
	n, err := r.DB.Read(bs)
	if err != nil {
		log.Fatalf("could not read from db: %v", err)
	}

	var reminders []models.Reminder
	err = json.Unmarshal(bs[:n], &reminders)
	if err != nil {
		log.Fatalf("could unrmashal json: %v", err)
	}

	remindersMap := map[int]models.Reminder{}
	originalOrder := map[int]int{}
	for i, reminder := range reminders {
		if filterFn == nil || filterFn(reminder) {
			remindersMap[reminder.ID] = reminder
			originalOrder[reminder.ID] = i
		}
	}
	return remindersMap, originalOrder
}

// NextID fetches the next DB AUTOINCREMENT id
func (r Reminder) NextID() int {
	return r.DB.GenerateID()
}

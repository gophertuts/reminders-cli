package models

import "time"

// Reminder represents the reminder data structure
type Reminder struct {
	ID         int           `json:"id"`
	Title      string        `json:"title"`
	Message    string        `json:"message"`
	Duration   time.Duration `json:"duration"`
	CreatedAt  time.Time     `json:"created_at"`
	ModifiedAt time.Time     `json:"modified_at"`
}

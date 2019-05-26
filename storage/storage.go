package storage

import "time"

type Storage interface {
	Data() []*Record
	Save(record Record) error
}

type Record struct {
	ID        string
	UserID    string
	Content   string
	CreatedAt time.Time
	RemindAt  time.Time
}

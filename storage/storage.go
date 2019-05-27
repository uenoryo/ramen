package storage

import (
	"fmt"
	"time"
)

type Storage interface {
	Load() error
	Data() []*Record
	Save(record *Record) error
}

type Record struct {
	ID        string
	UserID    string
	Content   string
	CreatedAt time.Time
	RemindAt  time.Time
}

func (r *Record) String() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s", r.ID, r.UserID, r.Content, r.CreatedAt.String(), r.RemindAt.String())
}

package storage

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const (
	timeFormat = "2006-01-02 15:04:05 MST"
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

func NewFromCSVLine(rows []string) (*Record, error) {
	if len(rows) != 5 {
		return nil, errors.Errorf("invalid csv line %v", rows)
	}

	createdAt, err := time.Parse(timeFormat, rows[3])
	if err != nil {
		return nil, errors.Wrapf(err, "error invalid created at data %s", rows[3])
	}

	remindAt, err := time.Parse(timeFormat, rows[4])
	if err != nil {
		return nil, errors.Wrapf(err, "error invalid remind at data %s", rows[4])
	}

	return &Record{rows[0], rows[1], rows[2], createdAt, remindAt}, nil
}

func (r *Record) String() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s", r.ID, r.UserID, r.Content, r.CreatedAt.Format(timeFormat), r.RemindAt.Format(timeFormat))
}

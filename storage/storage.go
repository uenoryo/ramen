package storage

import (
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
	Delete(id string) error
}

type Record struct {
	ID        string
	UserID    string
	Channel   string
	Content   string
	CreatedAt time.Time
	RemindAt  time.Time
}

func NewFromCSVLine(rows []string) (*Record, error) {
	if len(rows) != 6 {
		return nil, errors.Errorf("invalid csv line %v", rows)
	}

	createdAt, err := time.Parse(timeFormat, rows[4])
	if err != nil {
		return nil, errors.Wrapf(err, "error invalid created at data %s", rows[4])
	}

	remindAt, err := time.Parse(timeFormat, rows[5])
	if err != nil {
		return nil, errors.Wrapf(err, "error invalid remind at data %s", rows[5])
	}

	return &Record{rows[0], rows[1], rows[2], rows[3], createdAt, remindAt}, nil
}

func (r *Record) Strings() []string {
	return []string{r.ID, r.UserID, r.Channel, r.Content, r.CreatedAt.Format(timeFormat), r.RemindAt.Format(timeFormat)}
}

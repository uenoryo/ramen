package ramen

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

const schemaPath = "./ramen.sql"

type ramen struct {
	storage Storage
}

type Storage interface {
	init() error
	save(datetime, memo string) error
}

type database struct {
	dsn *sql.DB
}

func NewDatabase(dsn *sql.DB) *database {
	return &database{
		dsn: dsn,
	}
}

func (db *database) init() error {
	file, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return errors.Wrap(err, "error read sql file")
	}
	sql := strings.TrimSuffix(string(file), "\n")

	if err != nil {
		return errors.Wrap(err, "error connect db")
	}

	query := strings.Split(sql, ";")
	for _, q := range query {
		if q == "" {
			continue
		}
		if _, err = db.dsn.Exec(q); err != nil {
			return errors.Wrapf(err, "error exec query, [%s]", q)
		}
	}
	return nil
}

func (db *database) save(memo, datetime string) error {
	q := fmt.Sprintf("INSERT INTO reminder(`memo`, `remember_at`) VALUES('%s', '%s')", memo, datetime)
	_, err := db.dsn.Exec(q)
	return err
}

func NewRamen(storage Storage) *ramen {
	return &ramen{
		storage: storage,
	}
}

func (ramen *ramen) Init() error {
	return ramen.storage.init()
}

func (ramen *ramen) Set(memo, datetime string) error {
	err := ramen.storage.save(memo, datetime)
	if err != nil {
		return errors.Wrap(err, "error save reminder")
	}
	return nil
}

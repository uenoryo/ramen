package ramen

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const schemaPath = "./ramen.sql"

func Init() error {
	file, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return errors.Wrap(err, "error read sql file")
	}
	sql := strings.TrimSuffix(string(file), "\n")

	db, err := db()
	if err != nil {
		return errors.Wrap(err, "error connect db")
	}

	query := strings.Split(sql, ";")
	for _, q := range query {
		if q == "" {
			continue
		}
		if _, err = db.Exec(q); err != nil {
			return errors.Wrapf(err, "error exec query, [%s]", q)
		}
	}
	return nil
}

func db() (*sql.DB, error) {
	dbhost := os.Getenv("RAMEN_DB_HOST")
	dbname := os.Getenv("RAMEN_DB_NAME")
	user := os.Getenv("RAMEN_DB_USER")
	password := os.Getenv("RAMEN_DB_PASSWORD")
	if dbname == "" || user == "" {
		return nil, errors.New("error connect database, database name and user name are required")
	}
	return sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", user, password, dbhost, dbname))
}

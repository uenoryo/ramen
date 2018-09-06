package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uenoryo/ramen/ramen"
)

func main() {
	mysql, err := db()
	if err != nil {
		log.Printf("error initialize database, %s", err.Error())
		return
	}
	db := ramen.NewDatabase(mysql)
	ramen := ramen.NewRamen(db)
	if err := ramen.Init(); err != nil {
		log.Printf("error ramen init, %s", err.Error())
	}

	if err := ramen.Set("test", "2012/10/10"); err != nil {
		log.Printf("error ramen set, %s", err.Error())
	}
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

package main

import (
	"log"

	"database/sql"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "student"
)

// Insert a record into the database
func (d *Department) Insert(db *sql.DB) error {
    t := reflect.TypeOf(d).Elem()
    v := reflect.ValueOf(d).Elem()

    // Create the list of column names and placeholders
    columns := []string{}
    placeholders := []string{}
    values := []interface{}{}
    for i := 0; i < t.NumField(); i++ {
        columns = append(columns, t.Field(i).Name)
        placeholders = append(placeholders, "?")
        values = append(values, v.Field(i).Interface())
    }

    // Generate the INSERT statement
    query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
        t.Name(),
        strings.Join(columns, ", "),
        strings.Join(placeholders, ", "))

    // Execute the query
    _, err := db.Exec(query, values...)
    if err != nil {
        return err
    }
    return nil
}


type Department struct {
	ID       int
	DeptName string
	DeptCode string
}

func main() {
	// Connect to the database
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	department := Department{
		ID:       1,
		DeptName: "Computer Science",
		DeptCode: "CS",
	}
	if err := department.Insert(db); err != nil {
		log.Fatal(err)
	}
}

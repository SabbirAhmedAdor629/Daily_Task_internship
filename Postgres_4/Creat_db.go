package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to the "postgres" database.
	db, err := sql.Open("postgres", "postgres://user:password@localhost/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create the "mydatabase" database.
	_, err = db.Exec("CREATE DATABASE mydatabase")
	if err != nil {
		panic(err)
	}

	fmt.Println("Database created successfully!")
}

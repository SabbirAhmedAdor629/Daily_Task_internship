package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "sabbir123"
	dbname   = "cmddb"
)

func main() {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	// close database
	defer db.Close()

	//     // Delete
	// deleteStmt := `delete from "students" where std_id=$1`
	// _, e := db.Exec(deleteStmt, 6)
	// CheckError(e)

	// Update
	// updateStmt := `update "students" set "std_name"=$1 where "std_id"=$2`
	// _, e := db.Exec(updateStmt, "Mutakabbir Ahmed", 3)
	// CheckError(e)

	// Insertion
	//  // insert
	// insertStmt := `insert into "students"("std_id", "std_name", "std_program", "std_stream") values(4, 'Samin', 'School', 'Primary')`
	// _, e := db.Exec(insertStmt)
	// CheckError(e)
	// dynamic
	// insertDynStmt := `insert into "students"("std_id", "std_name", "std_program", "std_stream") values($1, $2, $3, $4)`
	// _, e := db.Exec(insertDynStmt, 6, "Mutashim", "Engineering", "Mechanica")
	// CheckError(e)

	// Checking Connection
	//     // check db
	// err = db.Ping()
	// CheckError(err)
	// fmt.Println("Connected!")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

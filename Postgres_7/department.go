package main

import (
	"database/sql"
	"fmt"
	"log"

	//	"log"
	"reflect"
	"strconv"
	"strings"

	//	"time"
	//	"log"
	//	"time"
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

type Department struct {
	ID       int
	DeptName string
	DeptCode string
}

func (d *Department) Insert(db *sql.DB)error {
	t := reflect.TypeOf(d).Elem()
	v := reflect.ValueOf(d).Elem()

	// Create the list of column names and placeholders
	columns := []string{}
	placeholders := []string{}
	values := []interface{}{}
	for i := 0; i < t.NumField(); i++ {
		columns = append(columns, t.Field(i).Name)
		placeholders = append(placeholders, "$"+strconv.Itoa(i+1))
		values = append(values, v.Field(i).Interface())
	}
	// fmt.Println(t.Name())
	// //fmt.Println(columns)
	// fmt.Println(strings.Join(columns, ", "))
	// //fmt.Println(placeholders)
	// fmt.Println(strings.Join(placeholders, ", "))

	//  INSERT 
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",t.Name(),strings.Join(columns, ", "),strings.Join(placeholders, ", "))

	// Execute the query
	_, err := db.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()

	department := &Department{
		ID:       87,
		DeptName: "SWE",
		DeptCode: "SWE-111214",
	}
	
	// department.Insert(db)

	if err := department.Insert(db); err != nil {
		log.Fatal(err)
	}

	
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

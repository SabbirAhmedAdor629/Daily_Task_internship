package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "student"
)

type DB struct {
	*sql.DB
}

func (db *DB) Connect() error {
	conn, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		return err
	}
	db.DB = conn
	return nil
}

func (db *DB) Insert(i interface{}) error {
	// interface type
	t := reflect.TypeOf(i).Elem()
	v := reflect.ValueOf(i).Elem()

	columns := []string{}
	placeholders := []string{}
	values := []interface{}{}
	for i := 0; i < t.NumField(); i++ {
		columns = append(columns, t.Field(i).Name)
		placeholders = append(placeholders, "$"+strconv.Itoa(i+1))
		values = append(values, v.Field(i).Interface())
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", t.Name(), strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	// Execute
	_, err := db.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

type Department struct {
	ID       int
	Dept_Name string
	DeptCode string
}

type Student struct {
	Id    int
	Name  string
	Email string
}


func main() {
	
	db := &DB{}
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//insert data into DEPARTMENT TABLE
	department := &Department{
		ID:       63654564,
		Dept_Name: "SWE",
		DeptCode: "",
	}
	if err := db.Insert(department); err != nil {
		log.Fatal(err)
	}

	//insert data into STUDENT TABLE
	student_1 := &Student{
		Id:    2554,
		Name:  "Ador",
		Email: "sdfasfaaf",
	}
	if err := db.Insert(student_1); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"database/sql"
	"fmt"

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


type Student struct {
	id    int
	name  string
	email string
}

func (s *Student) Insert(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO student(id, name, email) VALUES($1, $2, $3) RETURNING id",
		s.id, s.name, s.email,
	).Scan(&s.id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Student) Update(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE student SET name=$1, email=$2 WHERE id=$3",
		s.name, s.email, s.id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Student) Delete(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM student WHERE id=$1", s.id)
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

	student := &Student{
		id:    23,
		name:  "hasib",
		email: "hasib@gmail.com",
	}

	// Insert	Student
	student.Insert(db)

	// Update	Student
	// student.name = "Md. Alif"
	// student.email = "Md.Alif@gmail.com"
	// err = student.Update(db)
	// CheckError(err)

	// Delete	Student
	// err = student.Delete(db)
	// CheckError(err)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

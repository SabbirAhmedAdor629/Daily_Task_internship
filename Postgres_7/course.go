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


type Course struct {
	id          int
	CourseName  string
	CourseCode  string
}

func (c *Course) Insert(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO det(id, course_name, course_code) VALUES($1, $2, $3) RETURNING id",
		c.id, c.CourseName, c.CourseCode,
	).Scan(&c.id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Course) Update(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE course SET course_name=$1 WHERE id=$2",
		c.CourseName, c.id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *Course) Delete(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM Course WHERE id=$1", c.id)
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

	course := &Course{
	id:          55012,
	CourseName:  "DBMS",
	CourseCode:  "DB-1222",
	}

	err = course.Insert(db)
	CheckError(err)

	// // Update Course
	// course.CourseName = "srweww"
	// err = course.Update(db)
	// CheckError(err)

	// // Delete Course
	// err = course.Delete(db)
	// CheckError(err)


}

func CheckError(err error) {
	if err != nil {
	panic(err)
	}
}

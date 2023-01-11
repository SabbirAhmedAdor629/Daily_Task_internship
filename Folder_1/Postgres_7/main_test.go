package main

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "student"
)

func TestInsert(t *testing.T) {
	// Set up a connection to the database
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create a test case for the Insert function
	testCases := []struct {
		id          int
		CourseName  string
		CourseCode  string
		expectedErr error
	}{
		{
			id:          6,
			CourseName:  "DBMS",
			CourseCode:  "DB-1222",
			expectedErr: nil,
		},
		{
			id:          7,
			CourseName:  "",
			CourseCode:  "DB-1223",
			expectedErr: fmt.Errorf("Error inserting course: Course name cannot be empty"),
		},
	}

	// Iterate through the test cases
	for _, tc := range testCases {
		course := &Course{
			id:         tc.id,
			CourseName: tc.CourseName,
			CourseCode: tc.CourseCode,
		}
		err := course.Insert(db)
		if err != tc.expectedErr {
			t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
		}
	}
}

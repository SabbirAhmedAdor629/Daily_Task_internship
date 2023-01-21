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

func (db *DB) Update(i interface{}, id int) error {
    t := reflect.TypeOf(i).Elem()
    v := reflect.ValueOf(i).Elem()
    setStatements := []string{}
    values := []interface{}{}
    for i := 0; i < t.NumField(); i++ {
        setStatements = append(setStatements, t.Field(i).Name+"=$"+strconv.Itoa(i+1))
        values = append(values, v.Field(i).Interface())
    }
    values = append(values, id)
    
    query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", t.Name(), strings.Join(setStatements, ", "), len(values))
    
    // Execute
    _, err := db.Exec(query, values...)
    if err != nil {
        return err
    }
    return nil
}

func (db *DB) Delete(i interface{}, id int) error {
    t := reflect.TypeOf(i).Elem()
    query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", t.Name())

    // Execute
    _, err := db.Exec(query, id)
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

	//insert data into STUDENT TABLE
	// student_1 := &Student{
	// 	Id:    25905464,
	// 	Name:  "Ador",
	// 	Email: "sdfasfaaf",
	// }
	// if err := db.Insert(student_1); err != nil {
	// 	log.Fatal(err)
	// }

    // department_1:= &Department{
    //     ID:       87,
    //     Dept_Name: "SWE",
    //     DeptCode: "SWE-12dsr33",
    // }
    // // Call the Update function
    // if err := db.Update(department_1, 6075); err != nil {
    //     log.Fatal(err)
    // }

    //insert data into DEPARTMENT TABLE
	department := &Department{
		ID:       6075,
		Dept_Name: "SWE",
		DeptCode: "",
	}
	if err := db.Insert(department); err != nil {
		log.Fatal(err)
	}

    // update data in student table
    student_1 := &Student{
		Id:    1,
		Name:  "Ador",
		Email: "sdfasfaaf",
	}
	if err := db.Update(student_1,1); err != nil {
		log.Fatal(err)
	}

    // Delete from student table
    if err := db.Delete(&Student{}, 4); err != nil {
        log.Fatal(err)
    }

}

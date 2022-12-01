package main

import (
	"database/sql"
	"fmt"
//	"time"
	//	"log"
	//	"time"
	_ "github.com/lib/pq"
)
var db  *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "messaging_development"
)


func DeleteFromCampaign(db *sql.DB){
	var id int
	fmt.Scan(&id)
	del := fmt.Sprintf("DELETE FROM version_associations WHERE id = '%d'", id)
	_, err := db.Exec(del)
	CheckError(err)
}

// func UpdateCampaign(db *sql .DB) {
// 	var event_name string
// 	var id int
// 	fmt.Println("Enter new event name : ")
// 	fmt.Scan(&event_name)
// 	fmt.Println("Enter id ")
// 	fmt.Scan(&id)
// 	update	:=	`UPDATE version_associations SET "foreign_key_name" = $1 WHERE "id" = $2`
// 	_, err :=	db .Exec(update, event_name, id)
// 	CheckError(err)
// }


// func InsertIntoCampaign(db *sql.DB){
// 	//var id int
// 	var event_name string
// 	// created_at := time.Now().Format(time.RFC3339)
// 	// updated_at := time.Now().Format(time.RFC3339)
// 	// fmt.Println("Enter id : ")
// 	// fmt.Scan(&id)
// 	fmt.Println("Enter event name : ")
// 	fmt.Scan(&event_name)
// 	insert := fmt.Sprintf("INSERT INTO schema_migrations VALUES ('%s');", event_name )
// 	_, err := db.Exec(insert)
// 	CheckError(err)
// }


func main() {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	// close database
	defer db.Close()

	//InsertIntoCampaign(db)
	//UpdateCampaign(db)
	 DeleteFromCampaign(db)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

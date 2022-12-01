package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "messaging_development"
)

func main() {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	// close database
	defer db.Close()


						// Query
	// rows, err := db.Query(`SELECT "campaign_execution_log_id", "status" FROM "campaign_execution_log_chunks"`)
	// CheckError(err)
	
	// defer rows.Close()
	// for rows.Next() {
	// 	var campaign_execution_log_id int
	// 	var status string
	
	// 	err = rows.Scan(&campaign_execution_log_id, &status)
	// 	CheckError(err)
	
	// 	fmt.Println(campaign_execution_log_id, status)
	// }
	// CheckError(err)

	    			// Delete
	// deleteStmt := `delete from "students" where std_id=$1`
	// _, e := db.Exec(deleteStmt, 4)
	// CheckError(e)

						// Update
	// updateStmt := `UPDATE "campaign_execution_log_chunks" SET "campaign_id"=$1, "status"=$2 WHERE "id"=$3`
	// _, e := db.Exec(updateStmt, 514, "will be updated", 6)
	// CheckError(e)

	
						// Insertion
	//  // insert
	// insertStmt := `insert into "students"("std_id", "std_name", "std_program", "std_stream") values(4, 'Samin', 'School', 'Primary')`
	// _, e := db.Exec(insertStmt)
	// CheckError(e)

	// //dynamic
	// insertDynStmt := `insert into "campaign_execution_log_chunks"("id",
	// "campaign_execution_log_id", "campaign_id","utc_hour","execution_mins","processed_players",
	// "eligible_players","status","status_updated_at","distinct_players_sent","push_messages_sent",
	// "created_at","updated_at") values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	// _, e := db.Exec(insertDynStmt, 4, nil, 4, 8, 9, nil, 16, "not-updated", time.Now(), 18, 19, time.Now(), time.Now())
	// CheckError(e)

	// Checking Connection
	// err = db.Ping()
	// CheckError(err)
	// fmt.Println("Connected!")
}

				// CHECKING NULL VALUES
// func CheckNUllValue(key_1 int, tableName string){
// 	rows,err := db.Query(`SELECT "key_1" FROM "tableName"`);
// 	CheckError(err)
// 	for rows.Next() {
// 		var max_account *int
// 		err := rows.Scan(&max_account); 
// 		CheckError(err)
// 		if max_account == nil {
// 			fmt.Println("max_account is nil")
// 		}else{
// 			fmt.Println(*max_account);
// 		}
// 	}
// }

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

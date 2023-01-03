// package main

// import (
// 	"database/sql"
// 	"fmt"
// //	"time"
// //	"log"
// 	"time"
// 	_ "github.com/lib/pq"
// )
// var db  *sql.DB

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "postgres"
// 	dbname   = "messaging_development"
// )


// func ReadFromMessages(db *sql.DB){
// 	 // Execute a query to read the row with the given id.
// 	//  var event_name string
// 	//  err := db.QueryRow("SELECT event_name FROM campaigns WHERE id = $1", 1).Scan(&event_name)
// 	//  if err != nil {
// 	// 	 panic(err)
// 	//  }
// 	// fmt.Println(event_name)

// 	sqlStatement := `SELECT player_id FROM messages WHERE id=$1`
//     var Player_id int
//     err := db.QueryRow(sqlStatement, 1).Scan(&Player_id)
//     if err != nil {
//         panic(err)
//     }
//     fmt.Println(Player_id)
// }




// func InsertIntoMessages(db *sql.DB){
// 	// var Id int
// 	// var Name string
// 	// var Description string
// 	// var Parent_category_id int
// 	// var Priority int
// 	// var Can_opt_out bool
// 	// fmt.Println("Enter id : ")
// 	// fmt.Scan(&Id)
// 	// fmt.Println("Enter name : ")
// 	// fmt.Scan(&Name)
// 	// fmt.Println("Enter description : ")
// 	// fmt.Scan(&Description)
// 	// fmt.Println("Enter parent_catagory_id : ")
// 	// fmt.Scan(&Parent_category_id)
// 	// fmt.Println("Enter priority: ")
// 	// fmt.Scan(&Priority)	
// 	// fmt.Println("Enter can_opt_out: ")
// 	// fmt.Scan(&Can_opt_out)
// // 	email, loginTime := "human@example.com", time.Now()
// // result, err := db.Exec("INSERT INTO UserAccount VALUES ($1, $2)", email, loginTime)
// // if err != nil {
// //   panic(err)
// // }
// // insertTableAdId :=
// // 				`INSERT INTO ` + env.DbTable +
// // 					` (hashed_ad_id, app_ids, created_at, updated_at) VALUES ($1, $2, $3, $4) ` +
// // 					` ON CONFLICT DO NOTHING`

// 	// insert := fmt.Sprintf("INSERT INTO categories VALUES ('%s');", event_name )
// 	// _, err := db.Exec(insert)
	
// 	ViewedAt := time.Now()
// 	CreatedAt := time.Now()
// 	UpdatedAt := time.Now()
// 	CompletedAT := time.Now()
// 	CreatedDate := time.Now()
	

// 	_, err := db.Exec("INSERT INTO messages (id, player_id, message_template_id, viewed_at, created_at, updated_at, completed_at, created_date) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",1, 9, 10, ViewedAt, CreatedAt, UpdatedAt, CompletedAT, CreatedDate )
// 	CheckError(err)
// }



// func main() {
// 	// DB connection
// 	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
// 	db, err := sql.Open("postgres", psqlconn)
// 	CheckError(err)
// 	defer db.Close()

// 	//InsertIntoMessages(db)
// 	ReadFromMessages(db)
	
// }


// func CheckError(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Import the Postgres driver
)

// User represents a user in the database
type User struct {
	ID   int
	Name string
}

func main() {
		// DB connection
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()

	// Query the "users" table for a single row
	var user User
	err = db.QueryRow("SELECT * FROM users WHERE id = $1", 1).Scan(&user.ID, &user.Name)
	if err != nil {
		panic(err)
	}

	// Print the values of the struct
	fmt.Printf("ID: %d, Name: %s\n", user.ID, user.Name)
}

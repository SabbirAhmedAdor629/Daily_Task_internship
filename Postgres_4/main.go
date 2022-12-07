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
	password = "mypassword"
	dbname   = "mydatabase"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Test the connection to the database
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to PostgreSQL database!")

	// Create a table
	createTable := `
    CREATE TABLE IF NOT EXISTS products (
        product_id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        price NUMERIC(10,2) NOT NULL
    );
    `
	_, err = db.Exec(createTable)
	if err != nil {
		panic(err)
	}

	// Insert a row into the table
	insertRow := `
    INSERT INTO products (name, price)
    VALUES ('iPhone', 999.99);
    `
	_, err = db.Exec(insertRow)
	if err != nil {
		panic(err)
	}

	// Read a single row from the table
	var name string
	var price float64
	readRow := `
    SELECT name, price FROM products
    WHERE product_id = 1;
    `
	err = db.QueryRow(readRow).Scan(&name, &price)
	if err != nil {
		panic(err)
	}
	fmt.Println("Read row:", name, price)

	// Update a row in the table
	updateRow := `
    UPDATE products
    SET name = 'iPhone X', price = 1199.99
    WHERE product_id = 1;
    `
	_, err = db.Exec(updateRow)
	if err != nil {
		panic(err)
	}

	// Delete a row from the table
	deleteRow := `
    DELETE FROM products
    WHERE product_id = 1;
    `
	_, err = db.Exec(deleteRow)
	if err != nil {
		panic(err)
	}
}

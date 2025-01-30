package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func main() {
	// Ensure environment variables are set correctly
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Parse the flag for task to add
	task := flag.String("add", "GoGym", "adds task")
	flag.Parse()

	// Dummy ID (for simplicity)
	ID := 1

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO tasks (ID, task, status, Created) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(ID, *task, 0, time.Now())
	if err != nil {
		panic(err.Error())
	}

	// Success message
	fmt.Println("Task added successfully!")
}

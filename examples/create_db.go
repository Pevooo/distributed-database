package examples

import (
	"fmt"
)

func CreateAndPopulateDB() error {
	// Create database and tables
	createDBReq := DBRequest{
		Command:  "exec",
		Query:    "CREATE DATABASE IF NOT EXISTS mydb",
		IsMaster: true,
	}
	if _, err := SendQuery(createDBReq, "http://localhost:8080"); err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}

	// Create users table
	createTableReq := DBRequest{
		Command:  "exec",
		Query:    "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT UNIQUE, email TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)",
		IsMaster: true,
	}
	if _, err := SendQuery(createTableReq, "http://localhost:8080"); err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	// Insert some test data
	insertReq := DBRequest{
		Command: "exec",
		Query:   "INSERT INTO users (username, email) VALUES (?, ?)",
		Params:  []interface{}{"test_user", "test@example.com"},
	}
	if _, err := SendQuery(insertReq, "http://localhost:8080"); err != nil {
		return fmt.Errorf("failed to insert data: %v", err)
	}

	// Query the data
	selectReq := DBRequest{
		Command: "query",
		Query:   "SELECT * FROM users WHERE username = ?",
		Params:  []interface{}{"test_user"},
	}
	if resp, err := SendQuery(selectReq, "http://localhost:8080"); err != nil {
		return fmt.Errorf("failed to query data: %v", err)
	} else {
		fmt.Printf("✅ Retrieved data: %v\n", resp.Data)
	}

	fmt.Println("🎉 Database created and populated successfully!")
	return nil
}

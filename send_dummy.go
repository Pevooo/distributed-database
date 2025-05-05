package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DBRequest struct {
	Command      string        `json:"command"`
	Query        string        `json:"query,omitempty"`
	Params       []interface{} `json:"params,omitempty"`
	IsMaster     bool          `json:"is_master,omitempty"`
	IsReplicated bool          `json:"is_replicated,omitempty"`
}

func sendQuery(req DBRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	fmt.Printf("✅ Query executed successfully: %s\n", req.Query)
	return nil
}

func main() {
	// Create database and tables
	createDBReq := DBRequest{
		Command:  "exec",
		Query:    "CREATE DATABASE IF NOT EXISTS mydb",
		IsMaster: true,
	}
	if err := sendQuery(createDBReq); err != nil {
		fmt.Printf("❌ Failed to create database: %v\n", err)
		return
	}

	// Create users table
	createTableReq := DBRequest{
		Command:  "exec",
		Query:    "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT UNIQUE, email TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)",
		IsMaster: true,
	}
	if err := sendQuery(createTableReq); err != nil {
		fmt.Printf("❌ Failed to create table: %v\n", err)
		return
	}

	// Insert some test data
	insertReq := DBRequest{
		Command: "exec",
		Query:   "INSERT INTO users (username, email) VALUES (?, ?)",
		Params:  []interface{}{"test_user", "test@example.com"},
	}
	if err := sendQuery(insertReq); err != nil {
		fmt.Printf("❌ Failed to insert data: %v\n", err)
		return
	}

	// Query the data
	selectReq := DBRequest{
		Command: "query",
		Query:   "SELECT * FROM users WHERE username = ?",
		Params:  []interface{}{"test_user"},
	}
	if err := sendQuery(selectReq); err != nil {
		fmt.Printf("❌ Failed to query data: %v\n", err)
		return
	}

	fmt.Println("🎉 All operations completed successfully!")
}

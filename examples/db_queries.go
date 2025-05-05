package examples

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DBRequest struct {
	Command  string        `json:"command"`
	Query    string        `json:"query,omitempty"`
	Params   []interface{} `json:"params,omitempty"`
	IsMaster bool          `json:"is_master,omitempty"`
}

type DBResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SendQuery(req DBRequest, nodeURL string) (*DBResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(nodeURL+"/query", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse response
	var dbResp DBResponse
	if err := json.Unmarshal(respBody, &dbResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if !dbResp.Success {
		return &dbResp, fmt.Errorf("query failed: %s", dbResp.Error)
	}

	return &dbResp, nil
}

func RunExampleQueries() error {
	// Example 1: Create a new table (must be sent to master)
	createTableReq := DBRequest{
		Command:  "exec",
		Query:    "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT UNIQUE, email TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)",
		IsMaster: true,
	}
	if _, err := SendQuery(createTableReq, "http://localhost:8080"); err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}
	fmt.Println("✅ Table created successfully")

	// Example 2: Insert multiple users
	users := []struct {
		username string
		email    string
	}{
		{"john_doe", "john@example.com"},
		{"jane_smith", "jane@example.com"},
		{"bob_wilson", "bob@example.com"},
	}

	for _, user := range users {
		insertReq := DBRequest{
			Command: "exec",
			Query:   "INSERT INTO users (username, email) VALUES (?, ?)",
			Params:  []interface{}{user.username, user.email},
		}
		if _, err := SendQuery(insertReq, "http://localhost:8080"); err != nil {
			return fmt.Errorf("failed to insert user %s: %v", user.username, err)
		}
	}
	fmt.Println("✅ Users inserted successfully")

	// Example 3: Update user email
	updateReq := DBRequest{
		Command: "exec",
		Query:   "UPDATE users SET email = ? WHERE username = ?",
		Params:  []interface{}{"john.doe@newdomain.com", "john_doe"},
	}
	if _, err := SendQuery(updateReq, "http://localhost:8080"); err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	fmt.Println("✅ User updated successfully")

	// Example 4: Select users with pagination
	selectReq := DBRequest{
		Command: "query",
		Query:   "SELECT * FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?",
		Params:  []interface{}{2, 0}, // Get first 2 users
	}
	if resp, err := SendQuery(selectReq, "http://localhost:8080"); err != nil {
		return fmt.Errorf("failed to select users: %v", err)
	} else {
		fmt.Printf("✅ Retrieved users: %v\n", resp.Data)
	}

	// Example 5: Search users
	searchReq := DBRequest{
		Command: "query",
		Query:   "SELECT * FROM users WHERE username LIKE ? OR email LIKE ?",
		Params:  []interface{}{"%john%", "%example%"},
	}
	if resp, err := SendQuery(searchReq, "http://localhost:8080"); err != nil {
		return fmt.Errorf("failed to search users: %v", err)
	} else {
		fmt.Printf("✅ Search results: %v\n", resp.Data)
	}

	// Example 6: Get user statistics
	statsReq := DBRequest{
		Command: "query",
		Query:   "SELECT COUNT(*) as total_users, MAX(created_at) as latest_user FROM users",
	}
	if resp, err := SendQuery(statsReq, "http://localhost:8080"); err != nil {
		return fmt.Errorf("failed to get statistics: %v", err)
	} else {
		fmt.Printf("✅ User statistics: %v\n", resp.Data)
	}

	return nil
}

// Example of how to use the package
func Example() {
	if err := RunExampleQueries(); err != nil {
		fmt.Printf("❌ Error running queries: %v\n", err)
		return
	}
	fmt.Println("🎉 All queries executed successfully!")
}

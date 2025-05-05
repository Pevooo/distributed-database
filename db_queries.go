package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	Command  string        `json:"command"`
	Query    string        `json:"query,omitempty"`
	Params   []interface{} `json:"params,omitempty"`
	IsMaster bool          `json:"is_master,omitempty"`
}

func sendQuery(req Request, nodeURL string) {
	body, _ := json.Marshal(req)
	resp, err := http.Post(nodeURL+"/query", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("❌ Error sending query: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("✅ Query sent to %s\n", nodeURL)
}

func main() {
	// Example 1: Create a new table (must be sent to master)
	createTableReq := Request{
		Command:  "exec",
		Query:    "CREATE TABLE IF NOT EXISTS employees (id INTEGER PRIMARY KEY, name TEXT, age INTEGER, department TEXT)",
		IsMaster: true,
	}
	sendQuery(createTableReq, "http://localhost:8080")

	// Example 2: Insert data
	insertReq := Request{
		Command: "exec",
		Query:   "INSERT INTO employees (name, age, department) VALUES (?, ?, ?)",
		Params:  []interface{}{"John Doe", 30, "Engineering"},
	}
	sendQuery(insertReq, "http://localhost:8080")

	// Example 3: Update data
	updateReq := Request{
		Command: "exec",
		Query:   "UPDATE employees SET age = ? WHERE name = ?",
		Params:  []interface{}{31, "John Doe"},
	}
	sendQuery(updateReq, "http://localhost:8080")

	// Example 4: Select data
	selectReq := Request{
		Command: "query",
		Query:   "SELECT * FROM employees WHERE department = ?",
		Params:  []interface{}{"Engineering"},
	}
	sendQuery(selectReq, "http://localhost:8080")

	// Example 5: Delete data
	deleteReq := Request{
		Command: "exec",
		Query:   "DELETE FROM employees WHERE name = ?",
		Params:  []interface{}{"John Doe"},
	}
	sendQuery(deleteReq, "http://localhost:8080")

	// Example 6: Complex query with joins (if you have multiple tables)
	complexQuery := Request{
		Command: "query",
		Query:   "SELECT e.name, d.department_name FROM employees e JOIN departments d ON e.department = d.id WHERE e.age > ?",
		Params:  []interface{}{25},
	}
	sendQuery(complexQuery, "http://localhost:8080")

	// Example 7: Aggregate query
	aggregateQuery := Request{
		Command: "query",
		Query:   "SELECT department, COUNT(*) as employee_count, AVG(age) as avg_age FROM employees GROUP BY department",
	}
	sendQuery(aggregateQuery, "http://localhost:8080")
}

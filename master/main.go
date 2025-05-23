package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

type Request struct {
	Command      string        `json:"command"`
	Query        string        `json:"query,omitempty"`
	Params       []interface{} `json:"params,omitempty"`
	IsMaster     bool          `json:"is_master,omitempty"`
	IsReplicated bool          `json:"is_replicated,omitempty"`
}

var slaves []string
var db *sql.DB

func generateSlaves(n int, startPort int) []string {
	for i := 1; i <= n; i++ {
		addr := fmt.Sprintf("http://localhost:%d", startPort+i)
		slaves = append(slaves, addr)
	}
	return slaves
}

func replicate(req Request) {
	req.IsReplicated = true
    req.IsMaster =true
	body, _ := json.Marshal(req)
	for _, slave := range slaves {
		resp, err := http.Post(slave+"/query", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("❌ Failed to POST to %s: %v\n", slave, err)
			continue
		}
		defer resp.Body.Close()

		// Read response body
		resBody := new(bytes.Buffer)
		resBody.ReadFrom(resp.Body)

		fmt.Printf("✅ Response from %s [%s]: %s\n", slave, resp.Status, resBody.String())
	}
}

func isDatabaseOperation(query string) bool {
	query = strings.ToUpper(strings.TrimSpace(query))
	return strings.HasPrefix(query, "CREATE DATABASE") ||
		strings.HasPrefix(query, "DROP DATABASE") ||
		strings.HasPrefix(query, "CREATE TABLE")||
        strings.HasPrefix(query, "DROP TABLE")
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Error Received request:", req)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	} else {
		fmt.Println("Received request:", req)
	}

	// Handle regular queries
	switch req.Command {
	case "exec":
		_, err := db.Exec(req.Query, req.Params...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "success")
	case "query":
		rows, err := db.Query(req.Query, req.Params...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		columns, _ := rows.Columns()
		count := len(columns)
		results := []map[string]interface{}{}

		for rows.Next() {
			vals := make([]interface{}, count)
			ptrs := make([]interface{}, count)
			for i := range vals {
				ptrs[i] = &vals[i]
			}

			if err := rows.Scan(ptrs...); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			rowMap := make(map[string]interface{})
			for i, col := range columns {
				rowMap[col] = vals[i]
			}

			results = append(results, rowMap)
		}

		json.NewEncoder(w).Encode(results)
	default:
		http.Error(w, "Unsupported command", http.StatusBadRequest)
	}

	// Only replicate if it's not a database operation and not already replicated
	if !req.IsReplicated {
		replicate(req)
	}
}

func main() {
	n, err := strconv.Atoi(os.Args[1])
	if err != nil || n <= 0 {
		log.Fatal("Invalid database number")
	}

	generateSlaves(n, 8080)
	db, err = sql.Open("sqlite", "file:slave0.db")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/query", handleQuery)
	fmt.Println("Node listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

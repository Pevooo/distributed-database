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

	_ "github.com/mattn/go-sqlite3"
)

type Request struct {
	Command      string        `json:"command"`
	Query        string        `json:"query,omitempty"`
	Params       []interface{} `json:"params,omitempty"`
	IsMaster     bool          `json:"is_master,omitempty"`
	IsReplicated bool          `json:"is_replicated,omitempty"` // Flag to prevent re-replication
}

var db *sql.DB
var masterAddr string
var allNodes []string

func isDatabaseOperation(query string) bool {
	query = strings.ToUpper(strings.TrimSpace(query))
	return strings.HasPrefix(query, "CREATE DATABASE") ||
		strings.HasPrefix(query, "DROP DATABASE")
}

func replicate(req Request, sourceNode string) {
	// Set the replicated flag to prevent re-replication
	req.IsReplicated = true
	body, _ := json.Marshal(req)

	// Replicate to all nodes except the source
	for _, node := range allNodes {
		if node == sourceNode {
			continue // Skip the source node
		}

		resp, err := http.Post(node+"/query", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("❌ Failed to POST to %s: %v\n", node, err)
			continue
		}
		defer resp.Body.Close()

		// Read response body
		resBody := new(bytes.Buffer)
		resBody.ReadFrom(resp.Body)

		fmt.Printf("✅ Response from %s [%s]: %s\n", node, resp.Status, resBody.String())
	}
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

	// Check if it's a database operation
	if isDatabaseOperation(req.Query) {
		http.Error(w, "Slave nodes cannot perform database operations", http.StatusForbidden)
		return
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

	// Only replicate if this is not already a replicated request
	if !req.IsReplicated {
		// Get the source node's address
		sourceNode := r.RemoteAddr
		if sourceNode == "" {
			sourceNode = r.Host
		}

		// Replicate the operation to all other nodes
		replicate(req, sourceNode)
	}
}

func main() {
	var err error
	n, err := strconv.Atoi(os.Args[1])
	if err != nil || n <= 0 {
		log.Fatal("Invalid database number")
	}

	// Set up master address
	masterAddr = "http://localhost:8080"

	// Set up all node addresses
	allNodes = append(allNodes, masterAddr) // Add master
	for i := 1; i <= n; i++ {
		addr := fmt.Sprintf("http://localhost:%d", 8080+i)
		allNodes = append(allNodes, addr)
	}

	dbName := fmt.Sprintf("slave%d.db", n)
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/query", handleQuery)
	port := fmt.Sprintf(":%d", 8080+n)
	fmt.Printf("Node listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

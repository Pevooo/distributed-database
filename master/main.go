package main

import (
    "database/sql"
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "log"
    "strconv"
    "os"

    _ "github.com/mattn/go-sqlite3"
)

type Request struct {
    Command string            `json:"command"`
    Query   string            `json:"query,omitempty"`
    Params  []interface{}     `json:"params,omitempty"`
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

func handleQuery(w http.ResponseWriter, r *http.Request) {
    var req Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        fmt.Println("Error Received request:", req)
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    } else {
        fmt.Println("Received request:", req)
    }
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

    replicate(req) // Replicate the request to all slaves
}


func main() {
    n, err := strconv.Atoi(os.Args[1])
    if err != nil || n <= 0 {
        log.Fatal("Invalid database number")
    }

    generateSlaves(n, 8080)

    db, err = sql.Open("sqlite3", "slave0.db")
    if err != nil {
        log.Fatal(err)
    }


    http.HandleFunc("/query", handleQuery)
    fmt.Println("Node listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

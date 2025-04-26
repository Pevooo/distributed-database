package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "os"

    _ "github.com/mattn/go-sqlite3"
)

type Request struct {
    Command string            `json:"command"`
    Query   string            `json:"query,omitempty"`
    Params  []interface{}     `json:"params,omitempty"`
}

var db *sql.DB

func handleQuery(w http.ResponseWriter, r *http.Request) {
    var req Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        fmt.Println("Received request:", req)
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
}

func main() {
    var err error
    n, err := strconv.Atoi(os.Args[1])
    if err != nil || n <= 0 {
        log.Fatal("Invalid database number")
    }

    dbName := fmt.Sprintf("slave%d.db", n)
    db, err = sql.Open("sqlite3", dbName)
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/query", handleQuery)
    port := fmt.Sprintf(":%d", 8080 + n)
    fmt.Printf("Node listening on %s\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}

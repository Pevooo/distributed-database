package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "strconv"

    _ "github.com/mattn/go-sqlite3"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run init_slaves.go <number_of_slaves>")
        return
    }

    n, err := strconv.Atoi(os.Args[1])
    if err != nil || n <= 0 {
        log.Fatal("Invalid number of slaves")
    }

    for i := 0; i < n; i++ {
        dbName := fmt.Sprintf("slave%d.db", i)
        db, err := sql.Open("sqlite3", dbName)
        if err != nil {
            log.Fatalf("Failed to create DB %s: %v", dbName, err)
        }

        _, err = db.Exec(`
            CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY,
                name TEXT,
                email TEXT
            );
        `)
        if err != nil {
            log.Fatalf("Failed to create table in %s: %v", dbName, err)
        }

        fmt.Printf("Initialized %s\n", dbName)
        db.Close()
    }

    fmt.Println("All slave databases initialized.")
}

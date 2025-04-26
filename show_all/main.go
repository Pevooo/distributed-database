package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"os"
	
	_ "github.com/mattn/go-sqlite3"
)

func main() {
    var err error
    numSlaves, err := strconv.Atoi(os.Args[1])
    if err != nil || numSlaves <= 0 {
        log.Fatal("Invalid number")
    }              // change this to how many DBs you have
	tableName := "users"              // the table you want to query
	dbPathFormat := "slave%d.db"      // assumes files like slave0.db, slave1.db, ...

	for i := 0; i < numSlaves; i++ {
		dbPath := fmt.Sprintf(dbPathFormat, i)
		fmt.Printf("\n🔍 Querying %s\n", dbPath)

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Printf("❌ Failed to open %s: %v\n", dbPath, err)
			continue
		}

		rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
		if err != nil {
			log.Printf("❌ Query failed in %s: %v\n", dbPath, err)
			db.Close()
			continue
		}

		cols, _ := rows.Columns()
		fmt.Printf("📋 Columns: %v\n", cols)

		for rows.Next() {
			values := make([]interface{}, len(cols))
			pointers := make([]interface{}, len(cols))
			for i := range values {
				pointers[i] = &values[i]
			}

			err := rows.Scan(pointers...)
			if err != nil {
				log.Printf("❌ Row scan failed in %s: %v\n", dbPath, err)
				continue
			}

			record := make(map[string]interface{})
			for i, col := range cols {
				record[col] = values[i]
			}

			fmt.Printf("✅ Row: %v\n", record)
		}

		rows.Close()
		db.Close()
	}
}

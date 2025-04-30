package main

import (
    "bufio"
    "database/sql"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"

    _ "modernc.org/sqlite" 

)
var allowedTypes = map[string]bool{
    "INTEGER": true,
    "TEXT":    true,
    "REAL":    true,
    
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run init_slaves.go <number_of_slaves>")
        return
    }

    n, err := strconv.Atoi(os.Args[1])
    if err != nil || n <= 0 {
        log.Fatal("Invalid number of slaves")
    }

    reader := bufio.NewReader(os.Stdin)

    // Get table name
    fmt.Print("Enter table name: ")
    tableName, _ := reader.ReadString('\n')
    tableName = strings.TrimSpace(tableName)

    // Get number of attributes
    var attrCount int;
    for{
        fmt.Print("Enter number of attributes: ")
        attrCountStr, _ := reader.ReadString('\n')
        attrCountt, err := strconv.Atoi(strings.TrimSpace(attrCountStr))
        if err != nil || attrCountt <= 0 {
            fmt.Println("Invalid number of attributes")           
        }else{
            attrCount=attrCountt;
            break;
        }
    }
    

    columns := make([]string, attrCount+1)
    for i := 0; i <= attrCount; i++ {
        if i==0{
            columns[i] = fmt.Sprintf("%s %s PRIMARY KEY", "id", "INTEGER")
        }else{
            fmt.Printf("Enter name for attribute #%d: ", i)
            name, _ := reader.ReadString('\n')
            name = strings.TrimSpace(name)
            for{
                fmt.Printf("Enter type for attribute #%d (e.g., INTEGER, TEXT): ", i+1)
                typ, _ := reader.ReadString('\n')
                typ = strings.ToUpper(strings.TrimSpace(typ))
                if allowedTypes[typ] {
                    columns[i] = fmt.Sprintf("%s %s", name, typ)
                    break;
                }else{
                    fmt.Println("Invalid type. Allowed types: INTEGER, TEXT, REAL.")
                }
                
            }
            
        }
    }

    schema := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, strings.Join(columns, ", "))
    fmt.Println(schema);
    for i := 0; i < n; i++ {
        dbName := fmt.Sprintf("slave%d.db", i)
        db, err := sql.Open("sqlite", dbName)
        if err != nil {
            log.Fatalf("Failed to create DB %s: %v", dbName, err)
        }

        _, err = db.Exec(schema)
        if err != nil {
            log.Fatalf("Failed to create table in %s: %v", dbName, err)
        }

        fmt.Printf("✅ Initialized %s with table %s\n", dbName, tableName)
        db.Close()
    }

    fmt.Println("🎉 All slave databases initialized with custom schema.")
}

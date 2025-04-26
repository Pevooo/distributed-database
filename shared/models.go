package shared

type Record map[string]string

type Table struct {
    Name    string
    Columns []string
}

type Database struct {
    Name string
}

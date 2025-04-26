package main

import (
	"net/http"
	"bytes"
	"encoding/json"
)

type Request struct {
    Command string            `json:"command"`
    Query   string            `json:"query,omitempty"`
    Params  []interface{}     `json:"params,omitempty"`
}


func main() {

	req := Request{
		Command: "exec",
		Query:   "INSERT INTO users (id, name) VALUES (?, ?)",
		Params:  []interface{}{10, "test"},
	}
	body, _ := json.Marshal(req)

	http.Post("http://localhost:8080/query", "application/json", bytes.NewBuffer(body))
}
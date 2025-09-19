package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	server := NewSaver()
	server.Start()

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", register)

	fmt.Println("Server started on http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

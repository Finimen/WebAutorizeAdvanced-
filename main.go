package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func initDB() (*sql.DB, error) {
	var db, err = sql.Open("sqlite", "/.users.db")

	if err != nil {
		return nil, err
	}

	createTable := `
	CREATE TABLE IF NOT EXIST users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNQUE NOT NULL,
		password TEXT NOT NULL)`

	_, err = db.Exec(createTable)

	return db, err
}

func main() {
	server := NewSaver()
	server.Start()

	var db, err = initDB()

	if err != nil {
		log.Fatal("DB Initialize error")
		return
	}

	var userRepository = SQLRepository{
		bd: db,
	}

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", register)

	fmt.Println("Server started on http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

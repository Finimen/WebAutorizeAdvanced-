package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func initDB() (*sql.DB, error) {
	var db, err = sql.Open("sqlite", "./users.db")

	if err != nil {
		return nil, err
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL)`

	_, err = db.Exec(createTable)

	return db, err
}

func getKey() string {
	return "secretKey"
}

func startAutoriz(db *sql.DB) {
	var userRepository = SQLRepository{
		bd: db,
	}

	var hasher = BcryptHasher{}
	var key = getKey()

	var loginHandler = LoginHandler{
		Repo:   userRepository,
		Hasher: hasher,
		JwtKey: []byte(key),
	}

	var registerHandler = RegisterHandler{
		UserRepo: userRepository,
		Hasher:   hasher,
	}

	http.HandleFunc("/login", loginHandler.loginHandler)
	http.HandleFunc("/register", registerHandler.registerHandler)
	http.HandleFunc("/secret", middelwareHandler(secretHandler, key))
}

func main() {
	server := NewSaver()
	server.Start()

	var db, err = initDB()

	if err != nil {
		log.Fatal("DB Initialize error")
		return
	}

	startAutoriz(db)

	fmt.Println("Server started on http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

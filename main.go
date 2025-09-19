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

func startAutoriz(db *sql.DB) {
	var userRepository = SQLRepository{
		bd: db,
	}

	var hasher = BcryptHasher{}

	var loginHandler = LoginHandler{
		Repo:   userRepository,
		Hasher: hasher,
		JwtKey: []byte("secretKey"),
	}

	var registerHandler = RegisterHandler{
		UserRepo: userRepository,
		Hasher:   hasher,
	}

	http.HandleFunc("/login", loginHandler.loginHandler)
	http.HandleFunc("/register", registerHandler.registerHandler)
	http.HandleFunc("/secret", MiddelwareHandler)
}

func main() {
	server := NewSaver()
	server.Start()

	var db, err = initDB()

	startAutoriz(db)

	if err != nil {
		log.Fatal("DB Initialize error")
		return
	}

	fmt.Println("Server started on http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func initDB(config *Config) (*sql.DB, error) {
	var db, err = sql.Open("sqlite", config.DBPath)

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

func startAuth(db *sql.DB, limiter *RateLimiter) {
	var userRepository = SQLRepository{
		bd: db,
	}

	var hasher = BcryptHasher{}
	var key = getKey()

	var loginHandler = LoginHandler{
		Repo:   &userRepository,
		Hasher: &hasher,
		JwtKey: []byte(key),
	}

	var registerHandler = RegisterHandler{
		UserRepo: &userRepository,
		Hasher:   &hasher,
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/login", RateLimitMiddleware(limiter, loginHandler.loginHandler))
	http.HandleFunc("/register", RateLimitMiddleware(limiter, registerHandler.registerHandler))
	http.HandleFunc("/secret", middelwareHandler(secretHandler, key))
}

func main() {
	config := DefaultConfig()
	limiter := NewRateLimiter(config.RateLimit, config.RateWindow)

	saver := NewSaver()
	saver.Start()

	defer saver.Stop()

	var db, err = initDB(&config)

	if err != nil {
		log.Fatal("DB Initialize error")
		return
	}

	startAuth(db, limiter)

	fmt.Println("Server started on http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

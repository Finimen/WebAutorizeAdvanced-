package main

import "database/sql"

type SQLRepository struct {
	bd *sql.DB
}

func (r *SQLRepository) GetUserByUsername(name string) (string, error) {
	var password string
	row := r.bd.QueryRow("SELECT password FROM users WHERE username = ?", name)
	err := row.Scan(&password)
	return password, err
}

func (r *SQLRepository) CreateUser(name, hashedPassword string) error {
	_, err := r.bd.Exec("INSERT INTO users (username, password) VALUES (?, ?)", name, hashedPassword)
	return err
}

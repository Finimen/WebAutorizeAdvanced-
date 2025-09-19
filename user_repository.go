package main

import "database/sql"

type SQLRepository struct {
	bd *sql.DB
}

func GetUserByUsername(name string) (hashedPassword string, err error) {
	return "", nil
}

func CreateUser(name, hashedPassword string) error {
	return nil
}

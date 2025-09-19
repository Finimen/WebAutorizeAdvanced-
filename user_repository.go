package main

import "database/sql"

type SQLRepository struct {
	bd *sql.DB
}

func (r *SQLRepository) GetUserByUsername(name string) (hashedPassword string, err error) {
	return "", nil
}

func (r *SQLRepository) CreateUser(name, hashedPassword string) error {
	return nil
}

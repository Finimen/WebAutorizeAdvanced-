package main

import (
	"context"
	"database/sql"
)

type IRepository interface {
	GetUserByUsername(context.Context, string) (string, error)
	CreateUser(context.Context, string, string) error
}

type SQLRepository struct {
	bd *sql.DB
}

func (r *SQLRepository) GetUserByUsername(ctx context.Context, name string) (string, error) {
	var password string
	row := r.bd.QueryRowContext(ctx, "SELECT password FROM users WHERE username = ?", name)
	err := row.Scan(&password)
	return password, err
}

func (r *SQLRepository) CreateUser(ctx context.Context, name, hashedPassword string) error {
	_, err := r.bd.ExecContext(ctx, "INSERT INTO users (username, password) VALUES (?, ?)", name, hashedPassword)
	return err
}

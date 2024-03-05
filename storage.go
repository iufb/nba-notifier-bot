package main

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(context.Context, *Account) error
	DeleteAccount(int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	conSrt := os.Getenv("POSTGRES_URL")
	db, err := sql.Open("postgres", conSrt)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) CreateAccount(ctx context.Context, acc *Account) error {
	query := `
INSERT INTO accounts (telegram_id)
VALUES ($1);
    `
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, query, acc.TelegramId); err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `
    delete  from accounts where id=$1
    `
	_, err := s.db.Exec(query, id)
	return err
}

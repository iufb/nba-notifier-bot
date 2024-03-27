package main

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(context.Context, *Account) error
	GetAccounts(context.Context) ([]*Account, error)
	DeleteAccount(context.Context, int) error
	GetTeamUrl(context.Context, string) (string, error)
	AddTeamToFavourite(context.Context, int, string) error
	DeleteTeamFromFavourite(context.Context, int, string) error
	GetAccountFavouriteTeams(context.Context, int) ([]Team, error)
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

func (s *PostgresStore) GetAccounts(ctx context.Context) ([]*Account, error) {
	query := `
select * from accounts 
    `
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		acc := &Account{}
		err := rows.Scan(&acc.TelegramId)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (s *PostgresStore) DeleteAccount(ctx context.Context, id int) error {
	query := `
    delete  from accounts where telegram_id=$1
    `
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, query, id)
	return err
}

func (s *PostgresStore) GetTeamUrl(ctx context.Context, abbr string) (string, error) {
	query := `
    select url from teams where abbr=$1
    `
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	rows, err := conn.QueryContext(ctx, query, abbr)
	if err != nil {
		return "", err
	}

	var url string
	for rows.Next() {
		err := rows.Scan(&url)
		if err != nil {
			return "", err
		}
	}
	return url, nil
}

func (s *PostgresStore) AddTeamToFavourite(ctx context.Context, telegramId int, abbr string) error {
	query := `
INSERT INTO account_teams(telegram_id,team_abbr)
VALUES ($1,$2);
    `
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	_, err = conn.ExecContext(ctx, query, telegramId, abbr)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteTeamFromFavourite(ctx context.Context, telegramId int, teamName string) error {
	query := `
    delete from account_teams where  telegram_id=$1 and team_abbr=$2
    `
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, query, telegramId, teamName)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetAccountFavouriteTeams(ctx context.Context, telegramId int) ([]Team, error) {
	query := `
    select name, abbr ,url from account_teams join teams on account_teams.team_abbr = teams.abbr where telegram_id=$1`
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	rows, err := conn.QueryContext(ctx, query, telegramId)
	if err != nil {
		return nil, err
	}
	teams := []Team{}
	for rows.Next() {
		team := Team{}
		err := rows.Scan(&team.Name, &team.Abbr, &team.Url)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

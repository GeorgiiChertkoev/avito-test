package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"pr-reviewer/internal/config"
	"pr-reviewer/internal/domain/repo"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresUoW struct {
	db *sql.DB
}

func NewPostgresUoW(cfg config.Postgres) (*PostgresUoW, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DatabaseName,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Println("Failed to connect to postgres with dsn:", dsn)
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxConnections)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &PostgresUoW{db}, nil
}

func (u *PostgresUoW) Do(ctx context.Context, fn func(repos repo.Repositories) error) error {
	t, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	repos := repo.Repositories{
		TeamRepo: NewTeamRepo(t),
		UserRepo: NewUserRepo(t),
		PRRepo:   NewPRRepo(t),
	}

	if err := fn(repos); err != nil {
		_ = t.Rollback()
		return err
	}

	return t.Commit()
}

func (u *PostgresUoW) Close() {
	_ = u.db.Close()
}

package db

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/webitel/wlog"
)

var ErrDBNoExists = pgx.ErrNoRows

// WithTxFunc represents a function that will be executed within transaction.
type WithTxFunc func(ctx context.Context, tx *ConnectionTx) error

type Connection struct {
	log  *wlog.Logger
	pool *pgxpool.Pool
	psql sq.StatementBuilderType
}

type ConnectionTx struct {
	tx   pgx.Tx
	conn *Connection
}

func New(ctx context.Context, log *wlog.Logger, dsn string) (*Connection, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %v", err)
	}

	cfg.ConnConfig.Tracer = newTracer(log)
	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %v", err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %v", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &Connection{
		log:  log,
		pool: dbpool,
		psql: psql,
	}, nil
}

func (c *Connection) STDLib() *sql.DB {
	return stdlib.OpenDBFromPool(c.pool)
}

// WithTx executes a function within transaction.
func (c *Connection) WithTx(ctx context.Context, fn WithTxFunc) error {
	t, err := c.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	if err = fn(ctx, &ConnectionTx{tx: t, conn: c}); err != nil {
		if errRollback := t.Rollback(ctx); errRollback != nil {
			return fmt.Errorf("rollback tx: %w", err)
		}

		return fmt.Errorf("withTxFunc: %w", err)
	}

	if err = t.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

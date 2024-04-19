package db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (c *Connection) Domain(ctx context.Context, name string) (int64, error) {
	sql, args, err := c.psql.Select("id").From("fresh.domain").Where(sq.Eq{"name": name}).ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query: %v", err)
	}

	var id int64
	if err := c.pool.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return 0, fmt.Errorf("query row: %w", err)
	}

	return id, nil
}

func (c *Connection) CreateDomain(ctx context.Context, name string) (int64, error) {
	sql, args, err := c.psql.Insert("fresh.domain").SetMap(map[string]interface{}{"name": name}).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	if err := c.pool.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return 0, fmt.Errorf("query row: %w", err)
	}

	return id, nil
}

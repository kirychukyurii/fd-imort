package db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/kirychukyurii/fd-import/models"
)

func (c *Connection) Domain(ctx context.Context, d *models.Domain) (*models.Domain, error) {
	query := c.psql.Select("id", "name").From("fresh.domain")
	if d.ID != 0 {
		query = query.Where(sq.Eq{"id": d.ID})
	}

	if d.Name != "" {
		query = query.Where(sq.Eq{"name": d.Name})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %v", err)
	}

	var domain models.Domain
	if err := c.pool.QueryRow(ctx, sql, args...).Scan(&domain.ID, &domain.Name); err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &domain, nil
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

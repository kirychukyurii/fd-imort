package db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/kirychukyurii/fd-import/models"
)

func (c *Connection) Attachment(ctx context.Context, domain, id int64) (*models.Attachment, error) {
	sql, args, err := c.psql.Select("id", "name", "content_type", "file_size").From("fresh.attachment").
		Where(sq.Eq{"id": id, "domain_id": domain}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var attachment models.Attachment
	if err := c.pool.QueryRow(ctx, sql, args...).Scan(&attachment.ID, &attachment.Name, &attachment.ContentType, &attachment.FileSize); err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &attachment, nil
}

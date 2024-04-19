package models

import (
	"time"
)

// Attachment represents a row in the fresh.attachment table
type Attachment struct {
	RowID      int64     `json:"row_id" db:"row_id"`
	ImportedAt time.Time `json:"imported_at" db:"imported_at"`
	DomainID   int64     `json:"domain_id" db:"domain_id"`

	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	ContentType string    `json:"content_type" db:"content_type"`
	FileSize    int64     `json:"size" db:"file_size"`
	URL         string    `json:"attachment_url" db:"url"`
	ThumbURL    string    `json:"thumb_url" db:"thumb_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

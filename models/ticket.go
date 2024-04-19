package models

import (
	"time"
)

// Ticket represents a row in the fresh.ticket table
type Ticket struct {
	Raw []byte `json:"-" db:"-"`

	RowID      int64     `json:"row_id" db:"row_id"`
	ImportedAt time.Time `json:"imported_at" db:"imported_at"`
	AWSKey     string    `json:"aws_key" db:"aws_key"`
	DomainID   int64     `json:"domain_id" db:"domain_id"`

	RequesterName string `json:"requester_name" db:"requester_name"`

	ID                   int64           `json:"id,omitempty" db:"id"`
	Archived             bool            `json:"archived,omitempty" db:"archived"`
	Meta                 any             `json:"meta,omitempty" db:"meta"`
	Name                 string          `json:"name,omitempty" db:"name"`
	CCEmails             []string        `json:"cc_emails,omitempty" db:"cc_emails"`
	TicketCCEmails       []string        `json:"ticket_cc_emails,omitempty" db:"ticket_cc_emails"`
	CompanyID            int64           `json:"company_id,omitempty" db:"company_id"`
	CustomFields         any             `json:"custom_fields,omitempty" db:"custom_fields"`
	Deleted              bool            `json:"deleted,omitempty" db:"deleted"`
	Description          string          `json:"description,omitempty" db:"description"`
	DescriptionText      string          `json:"description_text,omitempty" db:"description_text"`
	DueBy                *time.Time      `json:"due_by,omitempty" db:"due_by"`
	Email                string          `json:"email,omitempty" db:"email"`
	EmailConfigID        int64           `json:"email_config_id,omitempty" db:"email_config_id"`
	FacebookID           string          `json:"facebook_id,omitempty" db:"facebook_id"`
	FrDueBy              *time.Time      `json:"fr_due_by,omitempty" db:"fr_due_by"`
	FrEscalated          bool            `json:"fr_escalated,omitempty" db:"fr_escalated"`
	NrDueBy              *time.Time      `json:"nr_due_by,omitempty" db:"nr_due_by"`
	NrEscalated          bool            `json:"nr_escalated,omitempty" db:"nr_escalated"`
	FwdEmails            []string        `json:"fwd_emails,omitempty" db:"fwd_emails"`
	GroupID              int64           `json:"group_id,omitempty" db:"group_id"`
	IsEscalated          bool            `json:"is_escalated,omitempty" db:"is_escalated"`
	Phone                string          `json:"phone,omitempty" db:"phone"`
	Priority             int64           `json:"priority,omitempty" db:"priority"`
	ProductID            int64           `json:"product_id,omitempty" db:"product_id"`
	ReplyCCEmails        []string        `json:"reply_cc_emails,omitempty" db:"reply_cc_emails"`
	RequesterID          int64           `json:"requester_id,omitempty" db:"requester_id"`
	ResponderID          int64           `json:"responder_id,omitempty" db:"responder_id"`
	Source               int64           `json:"source,omitempty" db:"source"`
	Spam                 bool            `json:"spam,omitempty" db:"spam"`
	Status               int64           `json:"status,omitempty" db:"status"`
	Subject              string          `json:"subject,omitempty" db:"subject"`
	Tags                 []string        `json:"tags,omitempty" db:"tags"`
	ToEmails             []string        `json:"to_emails,omitempty" db:"to_emails"`
	TwitterID            string          `json:"twitter_id,omitempty" db:"twitter_id"`
	Type                 string          `json:"type,omitempty" db:"type"`
	CreatedAt            *time.Time      `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt            *time.Time      `json:"updated_at,omitempty" db:"updated_at"`
	Attachments          []*Attachment   `json:"attachments,omitempty" db:"attachments"`
	Conversations        []*Conversation `json:"conversations,omitempty" db:"conversations"`
	AssociationType      int64           `json:"association_type,omitempty" db:"association_type"`
	SourceAdditionalInfo string          `json:"source_additional_info,omitempty" db:"source_additional_info"`
	SupportEmail         string          `json:"support_email,omitempty" db:"support_email"`
	FormID               int64           `json:"form_id,omitempty" db:"form_id"`
}

package models

import "time"

// Conversation represents a row in the fresh.conversation table
type Conversation struct {
	RowID      int64     `json:"row_id" db:"row_id"`
	ImportedAt time.Time `json:"imported_at" db:"imported_at"`

	Id                   int64         `json:"id" db:"id"`
	TicketID             int64         `json:"ticket_id" db:"ticket_id"`
	Body                 string        `json:"body" db:"body"`
	BodyText             string        `json:"body_text" db:"body_text"`
	Incoming             bool          `json:"incoming" db:"incoming"`
	ToEmails             []string      `json:"to_emails" db:"to_emails"`
	Category             int64         `json:"category" db:"category"`
	FromEmail            string        `json:"from_email" db:"from_email"`
	CCEmails             []string      `json:"cc_emails" db:"cc_emails"`
	BCCEmails            []string      `json:"bcc_emails" db:"bcc_emails"`
	Private              bool          `json:"private" db:"private"`
	Source               int64         `json:"source" db:"source"`
	SourceAdditionalInfo string        `json:"source_additional_info" db:"source_additional_info"`
	SupportEmail         string        `json:"support_email" db:"support_email"`
	CloudFiles           []any         `json:"cloud_files" db:"cloud_files"`
	AssociationType      int64         `json:"association_type" db:"association_type"`
	EmailFailureCount    int64         `json:"email_failure_count" db:"email_failure_count"`
	ThreadID             int64         `json:"thread_id" db:"thread_id"`
	ThreadMessageID      int64         `json:"thread_message_id" db:"thread_message_id"`
	AutoResponse         bool          `json:"auto_response" db:"auto_response"`
	AutomationID         int64         `json:"automation_id" db:"automation_id"`
	AutomationTypeID     int64         `json:"automation_type_id" db:"automation_type_id"`
	OutgoingFailures     []any         `json:"outgoing_failures" db:"outgoing_failures"`
	UserID               int64         `json:"user_id" db:"user_id"`
	LastEditedAt         time.Time     `json:"last_edited_at" db:"last_edited_at"`
	LastEditedUserID     int64         `json:"last_edited_user_id" db:"last_edited_user_id"`
	CreatedAt            time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at" db:"updated_at"`
	Attachments          []*Attachment `json:"attachments" db:"attachments"`
}

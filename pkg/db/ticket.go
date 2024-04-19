package db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/kirychukyurii/fd-import/models"
)

// Ticket retrieves a row_id from the `fresh.ticket_raw` table based on given domain ID and AWS key.
// It builds a query using `psql` and executes it using `pool.QueryRow`.
// If any error occurs during query execution, it returns the error.
func (c *Connection) Ticket(ctx context.Context, domain int64, key string) (bool, error) {
	query, args, err := c.psql.Select("row_id").From("fresh.ticket_raw").
		Where(sq.Eq{"domain_id": domain, "aws_key": key}).
		Limit(1).ToSql()
	if err != nil {
		return false, fmt.Errorf("build query: %v", err)
	}

	var ticketKey string
	if err := c.pool.QueryRow(ctx, query, args...).Scan(&ticketKey); err != nil {
		return false, fmt.Errorf("query row: %w", err)
	}

	return true, nil
}

// CreateTicket calls a function `fn` within a transaction on a ConnectionTx instance.
// It first creates conversations, then attachments, then the ticket itself,
// and finally the raw ticket. If any error occurs, it returns the error.
// If successful, it commits the transaction.
func (c *Connection) CreateTicket(ctx context.Context, ticket *models.Ticket) error {
	fn := func(ctx context.Context, tx *ConnectionTx) error {
		if err := tx.createConversations(ctx, ticket.DomainID, ticket.Conversations); err != nil {
			return fmt.Errorf("ticket [%d](%d): %v", ticket.ID, ticket.RequesterID, err)
		}

		if err := tx.createAttachments(ctx, ticket.DomainID, ticket.Attachments); err != nil {
			return fmt.Errorf("ticket [%d](%d): %v", ticket.ID, ticket.RequesterID, err)
		}

		if err := tx.createTicket(ctx, ticket); err != nil {
			return fmt.Errorf("ticket [%d](%d): %v", ticket.ID, ticket.RequesterID, err)
		}

		if err := tx.createRAWTicket(ctx, ticket.DomainID, ticket.AWSKey, ticket.ID, ticket.RequesterID, ticket.Raw); err != nil {
			return fmt.Errorf("raw [%d](%d): %v", ticket.ID, ticket.RequesterID, err)
		}

		return nil
	}

	if err := c.WithTx(ctx, fn); err != nil {
		return err
	}

	return nil
}

// createTicket inserts a new ticket record into the fresh.ticket table within transaction.
// The method takes the ctx context and ticket as arguments and inserts them into the table.
// If any error occurs during the process, it returns the error.
func (c *ConnectionTx) createTicket(ctx context.Context, ticket *models.Ticket) error {
	attIDs := make([]int64, 0)
	for _, attachment := range ticket.Attachments {
		attIDs = append(attIDs, attachment.ID)
	}

	values := map[string]interface{}{
		"aws_key": ticket.AWSKey, "id": ticket.ID, "archived": ticket.Archived, "meta": ticket.Meta, "name": ticket.Name,
		"cc_emails": ticket.CCEmails, "ticket_cc_emails": ticket.TicketCCEmails, "company_id": ticket.CompanyID,
		"custom_fields": ticket.CustomFields, "deleted": ticket.Deleted, "description": ticket.Description,
		"description_text": ticket.DescriptionText, "due_by": ticket.DueBy, "email": ticket.Email,
		"email_config_id": ticket.EmailConfigID, "facebook_id": ticket.FacebookID, "fr_due_by": ticket.FrDueBy,
		"fr_escalated": ticket.FrEscalated, "nr_due_by": ticket.NrDueBy, "nr_escalated": ticket.NrEscalated,
		"fwd_emails": ticket.FwdEmails, "group_id": ticket.GroupID, "is_escalated": ticket.IsEscalated,
		"phone": ticket.Phone, "priority": ticket.Priority, "product_id": ticket.ProductID,
		"reply_cc_emails": ticket.ReplyCCEmails, "requester_id": ticket.RequesterID, "responder_id": ticket.ResponderID,
		"source": ticket.Source, "spam": ticket.Spam, "status": ticket.Status, "subject": ticket.Subject,
		"tags": ticket.Tags, "to_emails": ticket.ToEmails, "twitter_id": ticket.TwitterID, "type": ticket.Type,
		"created_at": ticket.CreatedAt, "updated_at": ticket.UpdatedAt, "association_type": ticket.AssociationType,
		"source_additional_info": ticket.SourceAdditionalInfo, "support_email": ticket.SupportEmail,
		"form_id": ticket.FormID, "domain_id": ticket.DomainID, "requester_name": ticket.RequesterName,
		"attachment_ids": attIDs,
	}

	sql, args, err := c.conn.psql.Insert("fresh.ticket").SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build query: %v", err)
	}

	_, err = c.tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec query: %v", err)
	}

	return nil
}

// createRAWTicket inserts a new raw ticket json record into the fresh.ticket_raw table within transaction.
// The method takes the domain, key, id, requesterID, and ticket as arguments and inserts them into the table.
// If any error occurs during the process, it returns the error.
func (c *ConnectionTx) createRAWTicket(ctx context.Context, domain int64, key string, id, requesterID int64, ticket []byte) error {
	values := map[string]interface{}{
		"aws_key":      key,
		"ticket_id":    id,
		"requester_id": requesterID,
		"ticket":       ticket,
		"domain_id":    domain,
	}

	query, args, err := c.conn.psql.Insert("fresh.ticket_raw").SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build query: %v", err)
	}

	_, err = c.tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec query: %v", err)
	}

	return nil
}

// createConversations inserts multiple conversation records into the fresh.conversation table within transaction.
// It iterates over the conversations and calls createConversation method to insert each
// conversation record. If any error occurs during the process, it returns the error
// and stops further processing.
func (c *ConnectionTx) createConversations(ctx context.Context, domain int64, conversations []*models.Conversation) error {
	for i, cc := range conversations {
		if err := c.createConversation(ctx, domain, cc); err != nil {
			return fmt.Errorf("create conversation [%d](%d): %v", i, cc.Id, err)
		}

		if err := c.createAttachments(ctx, domain, cc.Attachments); err != nil {
			return fmt.Errorf("from conversation [%d](%d): %v", i, cc.Id, err)
		}
	}

	return nil
}

// createConversation inserts a new conversation record into the fresh.conversation table within transaction.
func (c *ConnectionTx) createConversation(ctx context.Context, domain int64, conversation *models.Conversation) error {
	attIDs := make([]int64, 0)
	for _, att := range conversation.Attachments {
		attIDs = append(attIDs, att.ID)
	}

	values := map[string]interface{}{
		"id":                     conversation.Id,
		"ticket_id":              conversation.TicketID,
		"body":                   conversation.Body,
		"body_text":              conversation.BodyText,
		"incoming":               conversation.Incoming,
		"to_emails":              conversation.ToEmails,
		"category":               conversation.Category,
		"from_email":             conversation.FromEmail,
		"cc_emails":              conversation.CCEmails,
		"bcc_emails":             conversation.BCCEmails,
		"private":                conversation.Private,
		"source":                 conversation.Source,
		"source_additional_info": conversation.SourceAdditionalInfo,
		"support_email":          conversation.SupportEmail,
		"cloud_files":            conversation.CloudFiles,
		"association_type":       conversation.AssociationType,
		"email_failure_count":    conversation.EmailFailureCount,
		"thread_id":              conversation.ThreadID,
		"thread_message_id":      conversation.ThreadMessageID,
		"auto_response":          conversation.AutoResponse,
		"automation_id":          conversation.AutomationID,
		"automation_type_id":     conversation.AutomationTypeID,
		"outgoing_failures":      conversation.OutgoingFailures,
		"user_id":                conversation.UserID,
		"last_edited_at":         conversation.LastEditedAt,
		"last_edited_user_id":    conversation.LastEditedUserID,
		"created_at":             conversation.CreatedAt,
		"updated_at":             conversation.UpdatedAt,
		"domain_id":              domain,
		"attachment_ids":         attIDs,
	}

	query, args, err := c.conn.psql.Insert("fresh.conversation").SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build query: %v", err)
	}

	if _, err = c.tx.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

// createAttachments iterates over the attachments and calls createAttachment method
// to insert each attachment record into the fresh.attachment table within transaction. If any error occurs
// during the process, it returns the error and stops further processing.
func (c *ConnectionTx) createAttachments(ctx context.Context, domain int64, attachments []*models.Attachment) error {
	for i, att := range attachments {
		if err := c.createAttachment(ctx, domain, att); err != nil {
			return fmt.Errorf("create attachment [%d](%d/%s): %v", i, att.ID, att.Name, err)
		}
	}

	return nil
}

// createAttachment inserts a new attachment record into the fresh.attachment table within transaction.
func (c *ConnectionTx) createAttachment(ctx context.Context, domain int64, attachment *models.Attachment) error {
	values := map[string]interface{}{
		"id":           attachment.ID,
		"name":         attachment.Name,
		"content_type": attachment.ContentType,
		"file_size":    attachment.FileSize,
		"url":          attachment.URL,
		"thumb_url":    attachment.ThumbURL,
		"created_at":   attachment.CreatedAt,
		"updated_at":   attachment.UpdatedAt,
		"domain_id":    domain,
	}

	query, args, err := c.conn.psql.Insert("fresh.attachment").SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build query: %v", err)
	}

	if _, err := c.tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("exec query: %v", err)
	}

	return nil
}

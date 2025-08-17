package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/abiewardani/go-messaging-system/internal/models"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) CreateMessage(message *models.Message) error {
	query := "INSERT INTO messages (tenant_id, content) VALUES ($1, $2)"
	_, err := r.db.Exec(query, message.TenantID, message.Content)
	return err
}

func (r *MessageRepository) GetMessagesByTenant(tenantID string) ([]models.Message, error) {
	query := "SELECT id, tenant_id, content FROM messages WHERE tenant_id = $1"
	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.TenantID, &message.Content); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// ListMessagesWithCursor implements cursor-based pagination for messages
func (r *MessageRepository) ListMessagesWithCursor(ctx context.Context, tenantID string, cursor string, limit int) ([]models.Message, string, error) {
	query := `
        SELECT id, tenant_id, payload, created_at 
        FROM messages 
        WHERE tenant_id = $1 
        AND ($2 = '' OR id > $2)
        ORDER BY id
        LIMIT $3
    `

	rows, err := r.db.QueryContext(ctx, query, tenantID, cursor, limit)
	if err != nil {
		return nil, "", fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.TenantID, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, "", fmt.Errorf("failed to scan message row: %w", err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, "", fmt.Errorf("error iterating message rows: %w", err)
	}

	return messages, cursor, nil
}

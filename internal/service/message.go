package service

import (
	"context"

	"github.com/abiewardani/go-messaging-system/internal/models"
	"github.com/abiewardani/go-messaging-system/internal/repository"
	"github.com/abiewardani/go-messaging-system/pkg/metrics"
)

type MessageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

// func (s *MessageService) CreateMessage(ctx context.Context, message *models.Message) error {
// 	return s.repo.Create(ctx, message)
// }

// func (s *MessageService) GetMessagesByTenant(ctx context.Context, tenantID string) ([]models.Message, error) {
// 	return s.repo.FindByTenant(ctx, tenantID)
// }

// func (s *MessageService) DeleteMessage(ctx context.Context, messageID string) error {
// 	return s.repo.Delete(ctx, messageID)
// }

// ListMessages lists messages for a tenant with pagination.
func (ms *MessageService) ListMessages(ctx context.Context, tenantID, cursor string, limit int) ([]models.Message, string, error) {
	// Validate input parameters
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}

	// Get messages from repository with pagination
	messages, _, err := ms.repo.ListMessagesWithCursor(ctx, tenantID, cursor, limit+1) // Fetch one extra for next cursor
	if err != nil {
		return nil, "", err
	}

	// If we got more items than the limit, set the next cursor
	var resultMessages []models.Message
	var resultCursor string

	if len(messages) > limit {
		resultMessages = messages[:limit]
		lastMessage := messages[limit-1]
		resultCursor = lastMessage.ID // Using message ID as cursor
	} else {
		resultMessages = messages
		resultCursor = "" // No more messages
	}

	// Update metrics
	metrics.MessageProcessed.WithLabelValues(tenantID, "listed").Inc()

	return resultMessages, resultCursor, nil
}

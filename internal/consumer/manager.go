package consumer

import (
	"fmt"
	"log"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

// MessageHandler defines the interface for processing messages
type MessageHandler interface {
	ProcessMessage([]byte) error
}

// TenantManager manages tenant consumers
type TenantManager struct {
	mu       sync.Mutex
	tenants  map[string]*TenantConsumer
	amqpConn *amqp091.Connection
}

// TenantConsumer represents a consumer for a specific tenant
type TenantConsumer struct {
	TenantID    string
	Channel     *amqp091.Channel
	Queue       string
	StopChan    chan struct{}
	WorkerCount int32
	handler     MessageHandler
}

// NewTenantManager creates a new tenant manager
func NewTenantManager(conn *amqp091.Connection) (*TenantManager, error) {
	tm := &TenantManager{
		tenants:  make(map[string]*TenantConsumer),
		amqpConn: conn,
	}

	// Start connection monitoring
	tm.monitorConnection()

	return tm, nil
}

// AddTenant implementation with proper error handling
func (tm *TenantManager) AddTenant(tenantID string, workerCount int32, handler MessageHandler) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Check if tenant already exists
	if _, exists := tm.tenants[tenantID]; exists {
		return fmt.Errorf("tenant %s already exists", tenantID)
	}

	// Create channel for tenant
	ch, err := tm.amqpConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	// Set QoS
	if err := ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		ch.Close()
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Declare queue with dead letter exchange
	args := amqp091.Table{
		"x-dead-letter-exchange":    "dlx",
		"x-dead-letter-routing-key": fmt.Sprintf("dl.%s", tenantID),
	}

	queueName := fmt.Sprintf("tenant_%s_queue", tenantID)
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments
	)
	if err != nil {
		ch.Close()
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	consumer := &TenantConsumer{
		TenantID:    tenantID,
		Channel:     ch,
		Queue:       q.Name,
		StopChan:    make(chan struct{}),
		WorkerCount: workerCount,
		handler:     handler,
	}

	// Start consumer workers
	if err := consumer.startWorkers(); err != nil {
		ch.Close()
		return fmt.Errorf("failed to start workers: %w", err)
	}

	tm.tenants[tenantID] = consumer
	return nil
}

func (tc *TenantConsumer) startWorkers() error {
	for i := int32(0); i < tc.WorkerCount; i++ {
		msgs, err := tc.Channel.Consume(
			tc.Queue, // queue
			fmt.Sprintf("%s-worker-%d", tc.TenantID, i), // consumer
			false, // auto-ack
			false, // exclusive
			false, // no-local
			false, // no-wait
			nil,   // args
		)
		if err != nil {
			return fmt.Errorf("failed to start consumer: %w", err)
		}

		go func() {
			for {
				select {
				case msg, ok := <-msgs:
					if !ok {
						return
					}
					if err := tc.handler.ProcessMessage(msg.Body); err != nil {
						log.Printf("Failed to process message for tenant %s: %v", tc.TenantID, err)
						msg.Nack(false, true) // Requeue the message
						continue
					}
					msg.Ack(false)
				case <-tc.StopChan:
					return
				}
			}
		}()
	}
	return nil
}

// RemoveTenant removes a tenant consumer with cleanup
func (tm *TenantManager) RemoveTenant(tenantID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	consumer, exists := tm.tenants[tenantID]
	if !exists {
		return fmt.Errorf("tenant %s not found", tenantID)
	}

	// Signal workers to stop
	close(consumer.StopChan)

	// Delete queue
	if _, err := consumer.Channel.QueueDelete(
		consumer.Queue, // queue name
		false,          // ifUnused
		false,          // ifEmpty
		false,          // noWait
	); err != nil {
		return fmt.Errorf("failed to delete queue: %w", err)
	}

	// Close channel
	if err := consumer.Channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}

	delete(tm.tenants, tenantID)
	return nil
}

// GetTenant retrieves a tenant consumer
func (tm *TenantManager) GetTenant(tenantID string) *TenantConsumer {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	return tm.tenants[tenantID]
}

// monitorConnection monitors the RabbitMQ connection and handles reconnection
func (tm *TenantManager) monitorConnection() {
	notifyClose := make(chan *amqp091.Error)
	tm.amqpConn.NotifyClose(notifyClose)

	go func() {
		for err := range notifyClose {
			log.Printf("RabbitMQ connection closed: %v", err)
			// TODO: Implement reconnection logic with exponential backoff
		}
	}()
}

// Close cleanly shuts down the TenantManager and all consumers
func (tm *TenantManager) Close() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for tenantID := range tm.tenants {
		if err := tm.RemoveTenant(tenantID); err != nil {
			log.Printf("Error removing tenant %s during shutdown: %v", tenantID, err)
		}
	}

	if err := tm.amqpConn.Close(); err != nil {
		return fmt.Errorf("failed to close AMQP connection: %w", err)
	}

	return nil
}

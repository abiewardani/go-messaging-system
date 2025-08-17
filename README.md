# Go Multi-Tenant Messaging System

## Overview
The `go-messaging-system` is a multi-tenant messaging application built using Go, RabbitMQ, and PostgreSQL. It is designed to support dynamic consumer management, partitioned data storage, and configurable concurrency, making it suitable for scalable messaging solutions.

## Features
- **Multi-Tenant Support**: Isolates data and consumers for different tenants
- **Dynamic Consumer Management**: Allows for the addition and removal of consumers at runtime
- **Partitioned Data Storage**: Utilizes PostgreSQL for efficient data storage and retrieval
- **Configurable Concurrency**: Supports adjustable concurrency levels for message processing
- **Cursor-based Pagination**: Efficient message retrieval with cursor pagination
- **Prometheus Metrics**: Built-in monitoring and metrics
- **Graceful Shutdown**: Proper cleanup of resources during shutdown
- **JWT Authentication**: Secure tenant isolation

## Project Structure
```
go-messaging-system
├── cmd
│   └── server
│       └── main.go          # Entry point of the application
├── internal
│   ├── app
│   │   └── server.go        # HTTP server setup and routing
│   ├── config
│   │   └── config.go        # Configuration management
│   ├── consumer
│   │   └── manager.go       # Tenant consumer management
│   ├── database
│   │   ├── migrations
│   │   │   └── schema.sql   # SQL schema for messages
│   │   └── postgres.go      # PostgreSQL connection management
│   ├── messaging
│   │   ├── publisher.go     # Message publishing logic
│   │   ├── consumer.go      # Message consuming logic
│   │   └── rabbitmq.go      # RabbitMQ connection management
│   ├── models
│   │   ├── message.go       # Message model structure
│   │   └── tenant.go        # Tenant model structure
│   ├── repository
│   │   ├── message.go       # Database operations for messages
│   │   └── tenant.go        # Database operations for tenants
│   └── service
│       ├── message.go       # Business logic for messages
│       └── tenant.go        # Business logic for tenants
├── pkg
│   └── metrics
│       └── metric.go        # Prometheus metrics definitions
├── config
│   └── config.json         # Application configuration
└── README.md              # Project documentation
```

## Prerequisites
- Go 1.21+
- Docker and Docker Compose
- Make (optional)

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/abiewardani/go-messaging-system.git
cd go-messaging-system
```

2. Start dependencies using Docker Compose:
```bash
docker-compose up -d
```

3. Run the application:
```bash
go run cmd/server/main.go
```

## Configuration

Configuration is managed via `config/config.json`:

```json
{
  "database_url": "postgres://user:password@localhost:5432/messaging_system?sslmode=disable",
  "rabbitmq_url": "amqp://user:password@localhost:5672/",
  "concurrency": 3
}
```

## API Documentation

### Create Tenant
```bash
curl -X POST http://localhost:8080/api/v1/tenants \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tenant One",
    "worker_count": 3
  }'
```

### Update Tenant Concurrency
```bash
curl -X PUT http://localhost:8080/api/v1/tenants/tenant123/config/concurrency \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "worker_count": 5
  }'
```

### List Messages (with pagination)
```bash
curl "http://localhost:8080/api/v1/messages?cursor=xyz&limit=10" \
  -H "Authorization: Bearer <your-token>"
```

### Delete Tenant
```bash
curl -X DELETE http://localhost:8080/api/v1/tenants/tenant123 \
  -H "Authorization: Bearer <your-token>"
```

## Docker Deployment

Start all services using Docker Compose:
```bash
# Build and start services
docker-compose up --build -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

## Monitoring

- Prometheus metrics: `http://localhost:8080/metrics`
- RabbitMQ management: `http://localhost:15672`
- Health check: `http://localhost:8080/health`

## Development

### Running Tests
```bash
# Run all tests
go test ./... -v

# Run integration tests
go test ./internal/tests -tags=integration
```

### Making Changes
1. Create a new branch
2. Make your changes
3. Run tests
4. Submit a pull request

## License
MIT License - see LICENSE file for
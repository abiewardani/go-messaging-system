go-messaging-system
├── cmd
│   └── server
│       └── main.go
├── internal
│   ├── app
│   │   └── server.go
│   ├── config
│   │   └── config.go
│   ├── consumer
│   │   └── manager.go
│   ├── database
│   │   ├── migrations
│   │   │   └── schema.sql
│   │   └── postgres.go
│   ├── messaging
│   │   ├── publisher.go
│   │   ├── consumer.go
│   │   └── rabbitmq.go
│   ├── models
│   │   ├── message.go
│   │   └── tenant.go
│   ├── repository
│   │   ├── message.go
│   │   └── tenant.go
│   └── service
│       ├── message.go
│       └── tenant.go
├── pkg
│   └── utils
│       └── common.go
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
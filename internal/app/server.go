package app

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/abiewardani/go-messaging-system/internal/consumer"
	"github.com/abiewardani/go-messaging-system/internal/models"
	"github.com/abiewardani/go-messaging-system/internal/service"
	"github.com/abiewardani/go-messaging-system/pkg/metrics"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Router         *mux.Router
	tenantManager  *consumer.TenantManager
	messageService *service.MessageService
}

// Request/Response structures
type CreateTenantRequest struct {
	Name        string `json:"name"`
	WorkerCount int32  `json:"worker_count"`
}

type UpdateConcurrencyRequest struct {
	WorkerCount int32 `json:"worker_count"`
}

type ListMessagesResponse struct {
	Messages   []models.Message `json:"messages"`
	NextCursor string           `json:"next_cursor"`
}

// NewServer creates and returns a new Server instance.
func NewServer(tm *consumer.TenantManager, ms *service.MessageService) *Server {
	s := &Server{
		Router:         mux.NewRouter(),
		tenantManager:  tm,
		messageService: ms,
	}

	// Add middleware
	s.Router.Use(s.authMiddleware)
	s.Router.Use(s.loggingMiddleware)

	// API routes
	api := s.Router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/tenants", s.CreateTenant).Methods("POST")
	api.HandleFunc("/tenants/{id}", s.DeleteTenant).Methods("DELETE")
	api.HandleFunc("/tenants/{id}/config/concurrency", s.UpdateConcurrency).Methods("PUT")
	api.HandleFunc("/messages", s.ListMessages).Methods("GET")

	// Monitoring
	s.Router.Handle("/metrics", promhttp.Handler())
	s.Router.HandleFunc("/health", s.healthCheck).Methods("GET")

	return s
}

func (s *Server) CreateTenant(w http.ResponseWriter, r *http.Request) {
	var req CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := s.tenantManager.AddTenant(req.Name, req.WorkerCount, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	// Get tenant ID from context (set by auth middleware)
	ctxTenantID := r.Context().Value("tenant_id").(string)
	if tenantID != ctxTenantID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if err := s.tenantManager.RemoveTenant(tenantID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) UpdateConcurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	var req UpdateConcurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate worker count
	if req.WorkerCount < 1 || req.WorkerCount > 10 {
		http.Error(w, "Worker count must be between 1 and 10", http.StatusBadRequest)
		return
	}

	tenant := s.tenantManager.GetTenant(tenantID)
	if tenant == nil {
		http.Error(w, "Tenant not found", http.StatusNotFound)
		return
	}

	// Update metrics before changing worker count
	metrics.WorkerCount.WithLabelValues(tenantID).Set(float64(req.WorkerCount))

	// TODO: Implement worker count update in TenantManager
	w.WriteHeader(http.StatusOK)
}

func (s *Server) ListMessages(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	limit := 10 // default limit
	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	tenantID := r.Context().Value("tenant_id").(string)

	messages, nextCursor, err := s.messageService.ListMessages(r.Context(), tenantID, cursor, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ListMessagesResponse{
		Messages:   messages,
		NextCursor: nextCursor,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	health := struct {
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
	}{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// loggingMiddleware logs incoming HTTP requests.
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		// Simple log to stdout, replace with your logger if needed
		method := r.Method
		path := r.URL.Path
		remote := r.RemoteAddr
		println("[LOG]", method, path, remote, duration.String())
	})
}

// authMiddleware is a simple authentication middleware that sets tenant_id in context.
// Replace this with your actual authentication logic as needed.
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID == "" {
			http.Error(w, "Missing tenant ID", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = contextWithTenantID(ctx, tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// contextWithTenantID sets the tenant_id in the context.
func contextWithTenantID(ctx context.Context, tenantID string) context.Context {
	type key string
	const tenantKey key = "tenant_id"
	return context.WithValue(ctx, tenantKey, tenantID)
}

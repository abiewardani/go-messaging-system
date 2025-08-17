package models

// Message represents a message in the messaging system.
type Message struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

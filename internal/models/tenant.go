package models

type Tenant struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Config      string `json:"config"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

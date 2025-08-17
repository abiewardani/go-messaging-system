package repository

import (
	"context"
	"database/sql"

	"github.com/abiewardani/go-messaging-system/internal/models"
)

type TenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(ctx context.Context, tenant *models.Tenant) error {
	query := "INSERT INTO tenants (name, description) VALUES ($1, $2) RETURNING id"
	return r.db.QueryRowContext(ctx, query, tenant.Name, tenant.Description).Scan(&tenant.ID)
}

func (r *TenantRepository) GetTenantByID(id int) (*models.Tenant, error) {
	tenant := &models.Tenant{}
	query := "SELECT id, name, config FROM tenants WHERE id = $1"
	err := r.db.QueryRow(query, id).Scan(&tenant.ID, &tenant.Name, &tenant.Config)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

func (r *TenantRepository) UpdateTenant(tenant *models.Tenant) error {
	query := "UPDATE tenants SET name = $1, config = $2 WHERE id = $3"
	_, err := r.db.Exec(query, tenant.Name, tenant.Config, tenant.ID)
	return err
}

func (r *TenantRepository) DeleteTenant(id int) error {
	query := "DELETE FROM tenants WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

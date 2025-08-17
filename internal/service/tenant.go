package service

import (
	"context"

	"github.com/abiewardani/go-messaging-system/internal/models"
	"github.com/abiewardani/go-messaging-system/internal/repository"
)

type TenantService struct {
	repo repository.TenantRepository
}

func NewTenantService(repo repository.TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	return s.repo.Create(ctx, tenant)
}

// func (s *TenantService) GetTenantByID(ctx context.Context, id string) (*models.Tenant, error) {
// 	return s.repo.FindByID(ctx, id)
// }

// func (s *TenantService) UpdateTenant(ctx context.Context, tenant *models.Tenant) error {
// 	return s.repo.Update(ctx, tenant)
// }

// func (s *TenantService) DeleteTenant(ctx context.Context, id string) error {
// 	return s.repo.Delete(ctx, id)
// }

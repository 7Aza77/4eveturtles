package usecase

import (
	"context"
	"goevent/internal/entity"
	"goevent/internal/repository"
)

type OrganizationUseCase interface {
	Create(ctx context.Context, org entity.Organization) (int64, error)
	GetByID(ctx context.Context, id int64) (entity.Organization, error)
	List(ctx context.Context) ([]entity.Organization, error)
	Update(ctx context.Context, org entity.Organization) error
	Delete(ctx context.Context, id int64) error
}

type Organization struct {
	repo repository.OrganizationRepository
}

func NewOrganization(repo repository.OrganizationRepository) *Organization {
	return &Organization{repo: repo}
}

func (u *Organization) Create(ctx context.Context, org entity.Organization) (int64, error) {
	return u.repo.Create(ctx, org)
}

func (u *Organization) GetByID(ctx context.Context, id int64) (entity.Organization, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *Organization) List(ctx context.Context) ([]entity.Organization, error) {
	return u.repo.List(ctx)
}

func (u *Organization) Update(ctx context.Context, org entity.Organization) error {
	return u.repo.Update(ctx, org)
}

func (u *Organization) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

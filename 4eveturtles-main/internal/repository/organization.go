package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"goevent/internal/entity"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org entity.Organization) (int64, error)
	GetByID(ctx context.Context, id int64) (entity.Organization, error)
	List(ctx context.Context) ([]entity.Organization, error)
	Update(ctx context.Context, org entity.Organization) error
	Delete(ctx context.Context, id int64) error
}

type OrganizationPostgres struct {
	db *sqlx.DB
}

func NewOrganizationPostgres(db *sqlx.DB) *OrganizationPostgres {
	return &OrganizationPostgres{db: db}
}

func (r *OrganizationPostgres) Create(ctx context.Context, org entity.Organization) (int64, error) {
	var id int64
	query := `INSERT INTO organizations (name, description, university_id, group_chat_link, owner_id)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(ctx, query,
		org.Name, org.Description, org.UniversityID, org.GroupChatLink, org.OwnerID,
	).Scan(&id)
	return id, err
}

func (r *OrganizationPostgres) GetByID(ctx context.Context, id int64) (entity.Organization, error) {
	var org entity.Organization
	query := "SELECT * FROM organizations WHERE id = $1"
	err := r.db.GetContext(ctx, &org, query, id)
	return org, err
}

func (r *OrganizationPostgres) List(ctx context.Context) ([]entity.Organization, error) {
	var orgs []entity.Organization
	query := "SELECT * FROM organizations ORDER BY id ASC"
	err := r.db.SelectContext(ctx, &orgs, query)
	return orgs, err
}

func (r *OrganizationPostgres) Update(ctx context.Context, org entity.Organization) error {
	query := `UPDATE organizations SET name=$1, description=$2, university_id=$3, group_chat_link=$4, updated_at=CURRENT_TIMESTAMP
			  WHERE id=$5`
	_, err := r.db.ExecContext(ctx, query, org.Name, org.Description, org.UniversityID, org.GroupChatLink, org.ID)
	return err
}

func (r *OrganizationPostgres) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM organizations WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

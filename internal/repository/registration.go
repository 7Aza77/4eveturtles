package repository

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

type RegistrationRepository interface {
	Register(ctx context.Context, userId, eventId int64) error
	Unregister(ctx context.Context, userId, eventId int64) error
	GetParticipantsCount(ctx context.Context, eventId int64) (int, error)
}

type RegistrationPostgres struct {
	db *sqlx.DB
}

func NewRegistrationPostgres(db *sqlx.DB) *RegistrationPostgres {
	return &RegistrationPostgres{db: db}
}

func (r *RegistrationPostgres) Register(ctx context.Context, userId, eventId int64) error {
	query := "INSERT INTO event_registrations (user_id, event_id) VALUES ($1, $2)"
	_, err := r.db.ExecContext(ctx, query, userId, eventId)
	return err
}

func (r *RegistrationPostgres) Unregister(ctx context.Context, userId, eventId int64) error {
	query := "DELETE FROM event_registrations WHERE user_id = $1 AND event_id = $2"
	res, err := r.db.ExecContext(ctx, query, userId, eventId)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("registration not found")
	}
	return nil
}

func (r *RegistrationPostgres) GetParticipantsCount(ctx context.Context, eventId int64) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM event_registrations WHERE event_id = $1"
	err := r.db.GetContext(ctx, &count, query, eventId)
	return count, err
}

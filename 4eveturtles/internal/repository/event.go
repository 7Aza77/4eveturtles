package repository

import (
	"context"
	"fmt"
	"goevent/internal/entity"
	"strings"

	"github.com/jmoiron/sqlx"
)

type EventRepository interface {
	Create(ctx context.Context, event entity.Event) (int64, error)
	GetByID(ctx context.Context, id int64) (entity.Event, error)
	List(ctx context.Context, filter EventFilter) ([]entity.Event, error)
	Update(ctx context.Context, event entity.Event) error
	Delete(ctx context.Context, id int64) error
}

type EventFilter struct {
	Limit    int
	Offset   int
	SortBy   string
	Order    string
	Title    string
	Location string
	FromDate string
	ToDate   string
}

var allowedSortColumns = map[string]bool{
	"id":       true,
	"title":    true,
	"date":     true,
	"location": true,
}

type EventPostgres struct {
	db *sqlx.DB
}

func NewEventPostgres(db *sqlx.DB) *EventPostgres {
	return &EventPostgres{db: db}
}

func (r *EventPostgres) Create(ctx context.Context, event entity.Event) (int64, error) {
	var id int64
	query := `INSERT INTO events (title, description, date, location, max_participants, creator_id)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.db.QueryRowContext(ctx, query,
		event.Title, event.Description, event.Date,
		event.Location, event.MaxParticipants, event.CreatorID,
	).Scan(&id)
	return id, err
}

func (r *EventPostgres) GetByID(ctx context.Context, id int64) (entity.Event, error) {
	var event entity.Event
	query := "SELECT * FROM events WHERE id = $1"
	err := r.db.GetContext(ctx, &event, query, id)
	return event, err
}

func (r *EventPostgres) List(ctx context.Context, filter EventFilter) ([]entity.Event, error) {
	var events []entity.Event
	var args []interface{}
	var conditions []string

	if filter.Title != "" {
		args = append(args, "%"+filter.Title+"%")
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", len(args)))
	}
	if filter.Location != "" {
		args = append(args, "%"+filter.Location+"%")
		conditions = append(conditions, fmt.Sprintf("location ILIKE $%d", len(args)))
	}
	if filter.FromDate != "" {
		args = append(args, filter.FromDate)
		conditions = append(conditions, fmt.Sprintf("date >= $%d", len(args)))
	}
	if filter.ToDate != "" {
		args = append(args, filter.ToDate)
		conditions = append(conditions, fmt.Sprintf("date <= $%d", len(args)))
	}

	query := "SELECT * FROM events"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	sortBy := "id"
	if allowedSortColumns[filter.SortBy] {
		sortBy = filter.SortBy
	}

	order := "ASC"
	if strings.ToUpper(filter.Order) == "DESC" {
		order = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)

	args = append(args, filter.Limit, filter.Offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	err := r.db.SelectContext(ctx, &events, query, args...)
	return events, err
}

func (r *EventPostgres) Update(ctx context.Context, event entity.Event) error {
	query := `UPDATE events SET title=$1, description=$2, date=$3, location=$4, max_participants=$5, updated_at=CURRENT_TIMESTAMP
			  WHERE id=$6`
	_, err := r.db.ExecContext(ctx, query,
		event.Title, event.Description, event.Date,
		event.Location, event.MaxParticipants, event.ID,
	)
	return err
}

func (r *EventPostgres) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM events WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

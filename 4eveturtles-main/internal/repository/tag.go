package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"goevent/internal/entity"
)

type TagRepository interface {
	Create(ctx context.Context, name string) (int64, error)
	List(ctx context.Context) ([]entity.Tag, error)
	AddTagsToEvent(ctx context.Context, eventId int64, tagIds []int64) error
	GetTagsByEventID(ctx context.Context, eventId int64) ([]entity.Tag, error)
	RemoveTagsFromEvent(ctx context.Context, eventId int64) error
}

type TagPostgres struct {
	db *sqlx.DB
}

func NewTagPostgres(db *sqlx.DB) *TagPostgres {
	return &TagPostgres{db: db}
}

func (r *TagPostgres) Create(ctx context.Context, name string) (int64, error) {
	var id int64
	query := `INSERT INTO tags (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name RETURNING id`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&id)
	return id, err
}

func (r *TagPostgres) List(ctx context.Context) ([]entity.Tag, error) {
	var tags []entity.Tag
	query := "SELECT * FROM tags ORDER BY name ASC"
	err := r.db.SelectContext(ctx, &tags, query)
	return tags, err
}

func (r *TagPostgres) AddTagsToEvent(ctx context.Context, eventId int64, tagIds []int64) error {
	for _, tagId := range tagIds {
		query := `INSERT INTO event_tags (event_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
		if _, err := r.db.ExecContext(ctx, query, eventId, tagId); err != nil {
			return err
		}
	}
	return nil
}

func (r *TagPostgres) GetTagsByEventID(ctx context.Context, eventId int64) ([]entity.Tag, error) {
	var tags []entity.Tag
	query := `SELECT t.id, t.name FROM tags t
			  INNER JOIN event_tags et ON t.id = et.tag_id
			  WHERE et.event_id = $1
			  ORDER BY t.name ASC`
	err := r.db.SelectContext(ctx, &tags, query, eventId)
	return tags, err
}

func (r *TagPostgres) RemoveTagsFromEvent(ctx context.Context, eventId int64) error {
	query := "DELETE FROM event_tags WHERE event_id = $1"
	_, err := r.db.ExecContext(ctx, query, eventId)
	return err
}

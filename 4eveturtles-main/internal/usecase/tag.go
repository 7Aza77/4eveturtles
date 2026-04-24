package usecase

import (
	"context"
	"goevent/internal/entity"
	"goevent/internal/repository"
)

type TagUseCase interface {
	Create(ctx context.Context, name string) (int64, error)
	List(ctx context.Context) ([]entity.Tag, error)
	AddTagsToEvent(ctx context.Context, eventId int64, tagIds []int64) error
	GetTagsByEventID(ctx context.Context, eventId int64) ([]entity.Tag, error)
	RemoveTagsFromEvent(ctx context.Context, eventId int64) error
}

type Tag struct {
	repo repository.TagRepository
}

func NewTag(repo repository.TagRepository) *Tag {
	return &Tag{repo: repo}
}

func (u *Tag) Create(ctx context.Context, name string) (int64, error) {
	return u.repo.Create(ctx, name)
}

func (u *Tag) List(ctx context.Context) ([]entity.Tag, error) {
	return u.repo.List(ctx)
}

func (u *Tag) AddTagsToEvent(ctx context.Context, eventId int64, tagIds []int64) error {
	return u.repo.AddTagsToEvent(ctx, eventId, tagIds)
}

func (u *Tag) GetTagsByEventID(ctx context.Context, eventId int64) ([]entity.Tag, error) {
	return u.repo.GetTagsByEventID(ctx, eventId)
}

func (u *Tag) RemoveTagsFromEvent(ctx context.Context, eventId int64) error {
	return u.repo.RemoveTagsFromEvent(ctx, eventId)
}

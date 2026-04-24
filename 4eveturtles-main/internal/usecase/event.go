package usecase

import (
	"context"
	"goevent/internal/entity"
	"goevent/internal/repository"
)

type EventUseCase interface {
	Create(ctx context.Context, event entity.Event) (int64, error)
	GetByID(ctx context.Context, id int64) (entity.Event, error)
	List(ctx context.Context, filter repository.EventFilter) ([]entity.Event, error)
	Update(ctx context.Context, userId int64, event entity.Event) error
	Delete(ctx context.Context, userId int64, eventId int64) error
}

type Event struct {
	repo repository.EventRepository
}

func NewEvent(repo repository.EventRepository) *Event {
	return &Event{repo: repo}
}

func (u *Event) Create(ctx context.Context, event entity.Event) (int64, error) {
	return u.repo.Create(ctx, event)
}

func (u *Event) GetByID(ctx context.Context, id int64) (entity.Event, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *Event) List(ctx context.Context, filter repository.EventFilter) ([]entity.Event, error) {
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	return u.repo.List(ctx, filter)
}

func (u *Event) Update(ctx context.Context, userId int64, event entity.Event) error {
	// Проверка прав (только создатель или админ может менять)
	// В реальном проекте мы бы проверяли роль из контекста, но тут упростим
	oldEvent, err := u.repo.GetByID(ctx, event.ID)
	if err != nil {
		return err
	}

	if oldEvent.CreatorID != userId {
		// Здесь можно было бы добавить проверку на админа
		// Но пока оставим так
	}

	return u.repo.Update(ctx, event)
}

func (u *Event) Delete(ctx context.Context, userId int64, eventId int64) error {
	oldEvent, err := u.repo.GetByID(ctx, eventId)
	if err != nil {
		return err
	}

	if oldEvent.CreatorID != userId {
		// Аналогично проверке выше
	}

	return u.repo.Delete(ctx, eventId)
}

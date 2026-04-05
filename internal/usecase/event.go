package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"goevent/internal/entity"
	"goevent/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
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
	rdb  *redis.Client
}

func NewEvent(repo repository.EventRepository, rdb *redis.Client) *Event {
	return &Event{repo: repo, rdb: rdb}
}

func (u *Event) Create(ctx context.Context, event entity.Event) (int64, error) {
	id, err := u.repo.Create(ctx, event)
	if err == nil && u.rdb != nil {
		u.rdb.Del(ctx, "events:list:*")
	}
	return id, err
}

func (u *Event) GetByID(ctx context.Context, id int64) (entity.Event, error) {
	key := fmt.Sprintf("event:%d", id)

	if u.rdb != nil {
		val, err := u.rdb.Get(ctx, key).Result()
		if err == nil {
			var event entity.Event
			if err := json.Unmarshal([]byte(val), &event); err == nil {
				return event, nil
			}
		}
	}

	event, err := u.repo.GetByID(ctx, id)
	if err == nil && u.rdb != nil {
		data, _ := json.Marshal(event)
		u.rdb.Set(ctx, key, data, time.Minute*10)
	}

	return event, err
}

func (u *Event) List(ctx context.Context, filter repository.EventFilter) ([]entity.Event, error) {
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	// Кэшируем только базовый список без сложных фильтров для простоты
	if u.rdb != nil && filter.Limit == 10 && filter.Offset == 0 {
		val, err := u.rdb.Get(ctx, "events:list:default").Result()
		if err == nil {
			var events []entity.Event
			if err := json.Unmarshal([]byte(val), &events); err == nil {
				return events, nil
			}
		}
	}

	events, err := u.repo.List(ctx, filter)
	if err == nil && u.rdb != nil && filter.Limit == 10 && filter.Offset == 0 {
		data, _ := json.Marshal(events)
		u.rdb.Set(ctx, "events:list:default", data, time.Minute*5)
	}

	return events, err
}

func (u *Event) Update(ctx context.Context, userId int64, event entity.Event) error {
	oldEvent, err := u.repo.GetByID(ctx, event.ID)
	if err != nil {
		return err
	}

	if oldEvent.CreatorID != userId {
		// return fmt.Errorf("access denied")
	}

	err = u.repo.Update(ctx, event)
	if err == nil && u.rdb != nil {
		u.rdb.Del(ctx, fmt.Sprintf("event:%d", event.ID), "events:list:*")
	}
	return err
}

func (u *Event) Delete(ctx context.Context, userId int64, eventId int64) error {
	oldEvent, err := u.repo.GetByID(ctx, eventId)
	if err != nil {
		return err
	}

	if oldEvent.CreatorID != userId {
		// return fmt.Errorf("access denied")
	}

	err = u.repo.Delete(ctx, eventId)
	if err == nil && u.rdb != nil {
		u.rdb.Del(ctx, fmt.Sprintf("event:%d", eventId), "events:list:*")
	}
	return err
}

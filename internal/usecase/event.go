package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goevent/internal/entity"
	"goevent/internal/metrics"
	"goevent/internal/repository"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

type EventUseCase interface {
	Create(ctx context.Context, event entity.Event) (int64, error)
	GetByID(ctx context.Context, id int64) (entity.Event, error)
	List(ctx context.Context, filter repository.EventFilter) ([]entity.Event, error)
	Update(ctx context.Context, userId int64, role string, event entity.Event) error
	Delete(ctx context.Context, userId int64, role string, eventId int64) error
}

type Event struct {
	repo repository.EventRepository
	rdb  *redis.Client
}

func NewEvent(repo repository.EventRepository, rdb *redis.Client) *Event {
	return &Event{repo: repo, rdb: rdb}
}

func (u *Event) invalidateListCache(ctx context.Context) {
	if u.rdb == nil {
		return
	}
	keys, _ := u.rdb.Keys(ctx, "events:list:*").Result()
	if len(keys) > 0 {
		u.rdb.Del(ctx, keys...)
	}
}

func (u *Event) Create(ctx context.Context, event entity.Event) (int64, error) {
	if event.MaxParticipants < 0 {
		return 0, errors.New("max_participants cannot be negative")
	}
	if event.Date.Before(time.Now()) {
		return 0, errors.New("event date must be in the future")
	}

	id, err := u.repo.Create(ctx, event)
	if err == nil {
		u.invalidateListCache(ctx)
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
				metrics.CacheHitsTotal.With(prometheus.Labels{"key_type": "event"}).Inc()
				return event, nil
			}
		} else if err == redis.Nil {
			metrics.CacheMissesTotal.With(prometheus.Labels{"key_type": "event"}).Inc()
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

	if u.rdb != nil && filter.Limit == 10 && filter.Offset == 0 &&
		filter.Title == "" && filter.Location == "" && filter.FromDate == "" && filter.ToDate == "" {
		val, err := u.rdb.Get(ctx, "events:list:default").Result()
		if err == nil {
			var events []entity.Event
			if err := json.Unmarshal([]byte(val), &events); err == nil {
				metrics.CacheHitsTotal.With(prometheus.Labels{"key_type": "events_list"}).Inc()
				return events, nil
			}
		} else if err == redis.Nil {
			metrics.CacheMissesTotal.With(prometheus.Labels{"key_type": "events_list"}).Inc()
		}
	}

	events, err := u.repo.List(ctx, filter)
	if err == nil && u.rdb != nil && filter.Limit == 10 && filter.Offset == 0 &&
		filter.Title == "" && filter.Location == "" && filter.FromDate == "" && filter.ToDate == "" {
		data, _ := json.Marshal(events)
		u.rdb.Set(ctx, "events:list:default", data, time.Minute*5)
	}

	return events, err
}

func (u *Event) Update(ctx context.Context, userId int64, role string, event entity.Event) error {
	oldEvent, err := u.repo.GetByID(ctx, event.ID)
	if err != nil {
		return errors.New("event not found")
	}

	if oldEvent.CreatorID != userId && role != string(entity.RoleAdmin) && role != string(entity.RoleModerator) {
		return errors.New("access denied: only the creator can update the event")
	}

	err = u.repo.Update(ctx, event)
	if err == nil && u.rdb != nil {
		u.rdb.Del(ctx, fmt.Sprintf("event:%d", event.ID))
		u.invalidateListCache(ctx)
	}
	return err
}

func (u *Event) Delete(ctx context.Context, userId int64, role string, eventId int64) error {
	oldEvent, err := u.repo.GetByID(ctx, eventId)
	if err != nil {
		return errors.New("event not found")
	}

	if oldEvent.CreatorID != userId && role != string(entity.RoleAdmin) && role != string(entity.RoleModerator) {
		return errors.New("access denied: only the creator can delete the event")
	}

	err = u.repo.Delete(ctx, eventId)
	if err == nil && u.rdb != nil {
		u.rdb.Del(ctx, fmt.Sprintf("event:%d", eventId))
		u.invalidateListCache(ctx)
	}
	return err
}

package usecase

import (
	"context"
	"errors"
	"goevent/internal/repository"
	"time"
)

type RegistrationUseCase interface {
	Register(ctx context.Context, userId, eventId int64) error
	Cancel(ctx context.Context, userId, eventId int64) error
}

type Registration struct {
	repo      repository.RegistrationRepository
	eventRepo repository.EventRepository
}

func NewRegistration(repo repository.RegistrationRepository, eventRepo repository.EventRepository) *Registration {
	return &Registration{
		repo:      repo,
		eventRepo: eventRepo,
	}
}

func (u *Registration) Register(ctx context.Context, userId, eventId int64) error {
	event, err := u.eventRepo.GetByID(ctx, eventId)
	if err != nil {
		return errors.New("event not found")
	}

	if event.MaxParticipants > 0 {
		count, err := u.repo.GetParticipantsCount(ctx, eventId)
		if err != nil {
			return err
		}
		if count >= event.MaxParticipants {
			return errors.New("event is full")
		}
	}

	return u.repo.Register(ctx, userId, eventId)
}

func (u *Registration) Cancel(ctx context.Context, userId, eventId int64) error {
	event, err := u.eventRepo.GetByID(ctx, eventId)
	if err != nil {
		return errors.New("event not found")
	}

	if event.Date.Before(time.Now()) {
		return errors.New("cannot unregister from a past event")
	}

	return u.repo.Unregister(ctx, userId, eventId)
}

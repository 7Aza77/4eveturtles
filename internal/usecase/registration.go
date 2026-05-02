package usecase

import (
	"context"
	"errors"
	"goevent/internal/repository"
)

type RegistrationUseCase interface {
	Register(ctx context.Context, userId, eventId int64) error
	Cancel(ctx context.Context, userId, eventId int64) error
	GetParticipants(ctx context.Context, eventId int64) ([]int64, error)
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

	already, err := u.repo.IsRegistered(ctx, userId, eventId)
	if err != nil {
		return err
	}
	if already {
		return errors.New("you are already registered for this event")
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
	return u.repo.Unregister(ctx, userId, eventId)
}

func (u *Registration) GetParticipants(ctx context.Context, eventId int64) ([]int64, error) {
	_, err := u.eventRepo.GetByID(ctx, eventId)
	if err != nil {
		return nil, errors.New("event not found")
	}
	return u.repo.GetParticipants(ctx, eventId)
}

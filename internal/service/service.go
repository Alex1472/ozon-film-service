package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Alex1472/ozon-film-service/internal/model"
	"github.com/Alex1472/ozon-film-service/internal/repo"
)

type FilmService struct {
	filmRepo    filmRepo
	eventRepo   eventRepo
	eventSender eventSender
}

type filmRepo interface {
	List(ctx context.Context) ([]*model.Film, error)
	Describe(ctx context.Context, filmID uint64) (*model.Film, error)
	Create(ctx context.Context, name string, rating float64, shortDescription string) (uint64, error)
	Remove(ctx context.Context, filmID uint64) error
}

type eventRepo interface {
	AddCreated(ctx context.Context, film *model.Film) error
	AddRemoved(ctx context.Context, filmID uint64) error
	Lock(ctx context.Context, n uint64) ([]*model.FilmEvent, error)
	Unlock(ctx context.Context, ids []uint64) error
	Remove(ctx context.Context, ids []uint64) (bool, error)
}

type eventSender interface {
	SendCreated(film *model.Film) error
	SendUpdated(film *model.Film) error
	SendRemoved(filmId uint64) error
}

func NewFilmService(filmRepo filmRepo, eventRepo eventRepo, eventSender eventSender) *FilmService {
	return &FilmService{
		filmRepo:    filmRepo,
		eventRepo:   eventRepo,
		eventSender: eventSender,
	}
}

var (
	ErrInvalidID = errors.New("invalid id")
)

func (fs *FilmService) Describe(ctx context.Context, id uint64) (*model.Film, error) {
	const op = "FilmService.Describe"
	film, err := fs.filmRepo.Describe(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrFilmNotFound) {
			return nil, fmt.Errorf("%s %w", op, ErrInvalidID)
		}
		return nil, fmt.Errorf("%s %w", op, err)
	}

	//fs.eventRepo.Lock(1)
	//fs.eventRepo.Unlock(ctx, []uint64{1, 2})
	//fs.eventRepo.Remove(ctx, []uint64{1, 3})

	return film, nil
}

func (fs *FilmService) List(ctx context.Context) ([]*model.Film, error) {
	const op = "FilmService.List"
	films, err := fs.filmRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return films, nil
}

func (fs *FilmService) Create(ctx context.Context, name string, rating float64, shortDescription string) (uint64, error) {
	const op = "FilmService.Create"

	//TODO add transaction
	id, err := fs.filmRepo.Create(ctx, name, rating, shortDescription)
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	film := &model.Film{
		ID:               id,
		Name:             name,
		Rating:           rating,
		ShortDescription: shortDescription,
	}
	err = fs.eventRepo.AddCreated(ctx, film)
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	err = fs.eventSender.SendCreated(film)
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	return id, nil
}

func (fs *FilmService) Remove(ctx context.Context, filmID uint64) error {
	const op = "FilmService.Remove"

	//TODO add transaction
	err := fs.filmRepo.Remove(ctx, filmID)
	if err != nil {
		if errors.Is(err, repo.ErrFilmNotFound) {
			return fmt.Errorf("%s %w", op, ErrInvalidID)
		}
		return fmt.Errorf("%s %w", op, err)
	}

	err = fs.eventRepo.AddRemoved(ctx, filmID)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	err = fs.eventSender.SendRemoved(filmID)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	return nil
}

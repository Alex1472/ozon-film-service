package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Alex1472/ozon-film-service/internal/model"
	"github.com/Alex1472/ozon-film-service/internal/repo"
)

type FilmService struct {
	repo Repo
}

type Repo interface {
	List(ctx context.Context) ([]*model.Film, error)
	Describe(ctx context.Context, filmID uint64) (*model.Film, error)
	Create(ctx context.Context, name string, rating float64, shortDescription string) (uint64, error)
	Remove(ctx context.Context, filmID uint64) error
}

func NewFilmService(repo Repo) *FilmService {
	return &FilmService{
		repo: repo,
	}
}

var (
	ErrInvalidID = errors.New("invalid id")
)

func (fs *FilmService) Describe(ctx context.Context, id uint64) (*model.Film, error) {
	const op = "FilmService.Describe"
	film, err := fs.repo.Describe(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrFilmNotFound) {
			return nil, fmt.Errorf("%s %w", op, ErrInvalidID)
		}
		return nil, fmt.Errorf("%s %w", op, err)
	}
	return film, nil
}

func (fs *FilmService) List(ctx context.Context) ([]*model.Film, error) {
	const op = "FilmService.List"
	films, err := fs.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return films, nil
}

func (fs *FilmService) Create(ctx context.Context, name string, rating float64, shortDescription string) (uint64, error) {
	const op = "FilmService.Create"
	id, err := fs.repo.Create(ctx, name, rating, shortDescription)
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}
	return id, nil
}

func (fs *FilmService) Remove(ctx context.Context, filmID uint64) error {
	const op = "FilmService.Remove"
	err := fs.repo.Remove(ctx, filmID)
	if err != nil {
		if errors.Is(err, repo.ErrFilmNotFound) {
			return fmt.Errorf("%s %w", op, ErrInvalidID)
		}
		return fmt.Errorf("%s %w", op, err)
	}
	return nil
}

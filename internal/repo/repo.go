package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/Alex1472/ozon-film-service/internal/model"
	"slices"
)

//type repo struct {
//	db        *sqlx.DB
//	batchSize uint
//}
//
//// NewRepo returns Repo interface
//func NewRepo(db *sqlx.DB, batchSize uint) Repo {
//	return &repo{db: db, batchSize: batchSize}
//}
//
//func (r *repo) Describe(ctx context.Context, filmID uint64) (*model.Film, error) {
//	return nil, nil
//}

type repo struct {
	films  []*model.Film
	nextID uint64
}

func NewRepo() *repo {
	return &repo{nextID: 1}
}

var (
	ErrFilmNotFound = errors.New("film not found")
)

func (r *repo) List(_ context.Context) ([]*model.Film, error) {
	result := make([]*model.Film, len(r.films))
	copy(result, r.films)
	return result, nil
}

func (r *repo) Describe(_ context.Context, filmID uint64) (*model.Film, error) {
	const op = "repo.Describe"
	for _, v := range r.films {
		if v.ID == filmID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("%s %w", op, ErrFilmNotFound)
}

func (r *repo) Create(_ context.Context, name string, rating float64, shortDescription string) (id uint64, err error) {
	film := &model.Film{
		ID:               r.nextID,
		Name:             name,
		Rating:           rating,
		ShortDescription: shortDescription,
	}
	r.nextID++
	r.films = append(r.films, film)
	return film.ID, nil
}

func (r *repo) Remove(_ context.Context, id uint64) error {
	const op = "repo.Remove"
	idx := slices.IndexFunc(r.films, func(film *model.Film) bool {
		return film.ID == id
	})
	if idx == -1 {
		return fmt.Errorf("%s %w", op, ErrFilmNotFound)
	}
	r.films = slices.Delete(r.films, idx, idx+1)
	return nil
}

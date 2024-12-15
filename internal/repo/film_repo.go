package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Alex1472/ozon-film-service/internal/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

//type filmRepo struct {
//	db        *sqlx.DB
//	batchSize uint
//}
//
//// NewFilmRepo returns Repo interface
//func NewFilmRepo(db *sqlx.DB, batchSize uint) Repo {
//	return &filmRepo{db: db, batchSize: batchSize}
//}
//
//func (r *filmRepo) Describe(ctx context.Context, filmID uint64) (*model.Film, error) {
//	return nil, nil
//}

type filmRepo struct {
	db *sqlx.DB
}

func NewFilmRepo(db *sqlx.DB) *filmRepo {
	return &filmRepo{db: db}
}

var (
	ErrFilmNotFound = errors.New("film not found")
)

func (r *filmRepo) List(ctx context.Context) ([]*model.Film, error) {
	const op = "film_repo.List"

	qb := sq.Select("id", "name", "rating", "short_description").
		PlaceholderFormat(sq.Dollar).
		From("films")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	var films []*model.Film
	err = r.db.SelectContext(ctx, &films, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	return films, nil
}

func (r *filmRepo) Describe(ctx context.Context, filmID uint64) (*model.Film, error) {
	const op = "filmRepo.Describe"

	qb := sq.Select("id", "name", "rating", "short_description").
		PlaceholderFormat(sq.Dollar).
		From("films").
		Where(sq.Eq{"id": filmID})

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	var film model.Film
	err = r.db.GetContext(ctx, &film, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%s %w", op, ErrFilmNotFound)
	}

	return &film, nil
}

func (r *filmRepo) Create(ctx context.Context, name string, rating float64, shortDescription string) (id uint64, err error) {
	const op = "filmRepo.Describe"

	query := sq.Insert("films").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "rating", "short_description").
		Values(name, rating, shortDescription).
		Suffix("RETURNING id").
		RunWith(r.db)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}
	if !rows.Next() {
		return 0, fmt.Errorf("%s %w", op, sql.ErrNoRows)
	}

	err = rows.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	return id, nil
}

func (r *filmRepo) Remove(ctx context.Context, id uint64) error {
	const op = "filmRepo.Remove"

	query := sq.Delete("films").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		RunWith(r.db)

	result, err := query.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}
	if cnt == 0 {
		return ErrFilmNotFound
	}

	return nil
}

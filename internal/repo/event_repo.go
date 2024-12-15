package repo

import (
	"context"
	"fmt"
	"github.com/Alex1472/ozon-film-service/internal/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type eventRepo struct {
	db *sqlx.DB
}

func NewEventRepo(db *sqlx.DB) *eventRepo {
	return &eventRepo{db: db}
}

func (er *eventRepo) AddCreated(ctx context.Context, film *model.Film) error {
	const op = "event_repo.AddCreated"

	query := sq.Insert("films_events").
		PlaceholderFormat(sq.Dollar).
		Columns("film_id", "type", "status", "payload").
		Values(film.ID, model.Created, model.Deferred, film).
		RunWith(er.db)

	_, err := query.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}
	return nil
}

func (er *eventRepo) AddRemoved(ctx context.Context, filmID uint64) error {
	const op = "event_repo.AddRemoved"
	query := sq.Insert("films_events").
		PlaceholderFormat(sq.Dollar).
		Columns("film_id", "type", "status").
		Values(filmID, model.Removed, model.Deferred).
		RunWith(er.db)

	_, err := query.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}
	return nil
}

func (er *eventRepo) Lock(ctx context.Context, n uint64) ([]*model.FilmEvent, error) {
	const op = "event_repo"

	query := "WITH updated_rows AS " +
		"(SELECT id FROM films_events " +
		"WHERE status = $1 " +
		"LIMIT $2) " +
		"UPDATE films_events fe " +
		"SET status = $3 " +
		"FROM updated_rows ur " +
		"WHERE ur.id = fe.id " +
		"RETURNING fe.film_id, fe.type, fe.status, fe.payload"

	rows, err := er.db.QueryxContext(ctx, query, model.Deferred, n, model.Processed)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	var events []*model.FilmEvent
	for rows.Next() {
		var event model.FilmEvent
		err = rows.StructScan(&event)
		if err != nil {
			return nil, fmt.Errorf("%s %w", op, err)
		}
		events = append(events, &event)
	}
	return events, nil
}

func (er *eventRepo) Unlock(ctx context.Context, ids []uint64) error {
	const op = "eventRepo.Unlock"

	query := sq.Update("films_events").
		PlaceholderFormat(sq.Dollar).
		Set("status", model.Deferred).
		Where(sq.Eq{"id": ids}).
		RunWith(er.db)

	_, err := query.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}
	return nil
}

func (er *eventRepo) Remove(ctx context.Context, ids []uint64) (bool, error) {
	const op = "eventRepo.Remove"

	query := sq.Delete("films_events").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": ids}).
		RunWith(er.db)

	_, err := query.ExecContext(ctx)
	if err != nil {
		return false, fmt.Errorf("%s %w", op, err)
	}
	return true, nil
}

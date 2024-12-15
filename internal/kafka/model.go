package kafka

import "github.com/Alex1472/ozon-film-service/internal/model"

type EventType uint8

const (
	_ EventType = iota
	Created
	Updated
	Removed
)

type Film struct {
	ID               uint64  `db:"id"`
	Name             string  `db:"name"`
	Rating           float64 `db:"rating"`
	ShortDescription string  `db:"short_description"`
}

type FilmEvent struct {
	ID   uint64    `db:"film_id"`
	Type EventType `db:"type"`
	Film *Film     `db:"payload"`
}

func NewKafkaEvent(filmId uint64, eventType EventType, film *model.Film) *FilmEvent {
	var kafkaFilm *Film
	if film != nil {
		kafkaFilm = toKafkaFilm(film)
	}
	return &FilmEvent{
		ID:   filmId,
		Type: eventType,
		Film: kafkaFilm,
	}
}

func toKafkaFilm(film *model.Film) *Film {
	return &Film{
		ID:               film.ID,
		Name:             film.Name,
		Rating:           film.Rating,
		ShortDescription: film.ShortDescription,
	}
}

package model

type EventType uint8

type EventStatus uint8

const (
	_ EventType = iota
	Created
	Updated
	Removed
)

const (
	_ EventStatus = iota
	Deferred
	Processed
)

type FilmEvent struct {
	ID     int64       `db:"film_id"`
	Type   EventType   `db:"type"`
	Status EventStatus `db:"status"`
	Film   *Film       `db:"payload"`
}

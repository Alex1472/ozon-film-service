package model

type Film struct {
	ID               uint64  `db:"id"`
	Name             string  `db:"name"`
	Rating           float64 `db:"rating"`
	ShortDescription string  `db:"short_description"`
}

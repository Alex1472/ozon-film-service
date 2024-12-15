package model

import (
	"encoding/json"
	"errors"
)

type Film struct {
	ID               uint64  `db:"id"`
	Name             string  `db:"name"`
	Rating           float64 `db:"rating"`
	ShortDescription string  `db:"short_description"`
}

func (f *Film) Scan(src interface{}) (err error) {
	if src == nil {
		return nil
	}

	var film Film
	switch src := src.(type) {
	case string:
		err = json.Unmarshal([]byte(src), &film)
	case []byte:
		err = json.Unmarshal(src, &film)
	default:
		return errors.New("Incompatible type")
	}

	if err != nil {
		return err
	}

	*f = film
	return nil
}

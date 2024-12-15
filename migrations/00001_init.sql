-- +goose Up
CREATE TABLE films (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  rating FLOAT NOT NULL,
  short_description TEXT NOT NULL,
  removed BOOLEAN DEFAULT FALSE,
  created TIMESTAMP NOT NULL DEFAULT NOW(),
  updated TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE films_events (
  id BIGSERIAL PRIMARY KEY,
  film_id BIGINT,
  type INTEGER NOT NULL,
  status INTEGER NOT NULL,
  payload JSONB,
  updated TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE filmsEvents;
DROP TABLE films;

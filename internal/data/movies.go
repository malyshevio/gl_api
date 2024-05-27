package data

import "time"

type Movie struct {
	ID        int64
	CreatedAt time.Time // timestamp when added to DB
	Title     string
	Year      int32
	Runtime   int32
	Genres    []string
	Version   int32
}

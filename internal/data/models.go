package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("record edit conflict")
)

type Models struct {
	Movies      MovieModel
	Permissions PermissionModel
	Users       UserModel
	Tokens      TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:      MovieModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Users:       UserModel{DB: db},
		Tokens:      TokenModel{DB: db},
	}
}

// type Models struct {
// 	Movies interface {
// 		Insert(movie *Movie) error
// 		Get(id int64) (*Movie, error)
// 		Update(movie *Movie) error
// 		Delete(id int64) error
// 	}
// }

// func NewModels(db *sql.DB) Models {
// 	return Models{
// 		Movies: MockMovieModel{},
// 	}
// }

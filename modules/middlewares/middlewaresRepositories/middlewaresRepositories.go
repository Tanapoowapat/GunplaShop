package middlewaresRepositories

import "github.com/jmoiron/sqlx"

type IMiddlewaresRepositories interface {
}

type middlewaresRepositories struct {
	db *sqlx.DB
}

func NewMiddlewaresRepositories(db *sqlx.DB) IMiddlewaresRepositories {
	return &middlewaresRepositories{
		db: db,
	}
}

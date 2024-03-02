package middlewaresRepositories

import (
	"fmt"

	"github.com/Tanapoowapat/GunplaShop/modules/middlewares"
	"github.com/jmoiron/sqlx"
)

type IMiddlewaresRepositories interface {
	FindAcessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresRepositories struct {
	db *sqlx.DB
}

func NewMiddlewaresRepositories(db *sqlx.DB) IMiddlewaresRepositories {
	return &middlewaresRepositories{
		db: db,
	}
}

func (r *middlewaresRepositories) FindAcessToken(userId, accessToken string) bool {

	query := `
		SELECT
			(CASE WHEN COUNT(*) = 1 THEN true ELSE false END)
		FROM "oauth"
		WHERE "user_id" = $1 AND "access_tokens" = $2
	`
	var check bool
	if err := r.db.Get(&check, query, userId, accessToken); err != nil {
		return false
	}
	return true
}

func (r *middlewaresRepositories) FindRole() ([]*middlewares.Role, error) {
	query := `
			SELECT
					"id",
					"title"
			FROM "roles"
			ORDER BY "id" DESC
	`
	roles := make([]*middlewares.Role, 0)
	if err := r.db.Select(&roles, query); err != nil {
		return nil, fmt.Errorf("roles are empty")
	}
	return roles, nil
}

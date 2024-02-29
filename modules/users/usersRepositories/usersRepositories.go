package usersRepositories

import (
	"context"
	"fmt"
	"time"

	"github.com/Tanapoowapat/GunplaShop/modules/users"
	"github.com/Tanapoowapat/GunplaShop/modules/users/usersPatterns"
	"github.com/jmoiron/sqlx"
)

type IUserRepositories interface {
	InsertUser(req *users.UserRegisterRequest, isAdmin bool) (*users.UserPassport, error)
	FindUserByEmail(email string) (*users.UserCredentialsCheck, error)
	InsertOauth(req *users.UserPassport) error
	FindOneOauth(refreshToken string) (*users.OAuth, error)
	UpdateOauth(req *users.UserTokens) error
	GetProfile(userId string) (*users.User, error)
	DeleteOAuth(oauthId string) error
}

type userRepositories struct {
	db *sqlx.DB
}

func UsersRepositories(db *sqlx.DB) IUserRepositories {
	return &userRepositories{
		db: db,
	}
}

func (r *userRepositories) InsertUser(req *users.UserRegisterRequest, isAdmin bool) (*users.UserPassport, error) {
	//
	result := usersPatterns.InsertUser(r.db, req, isAdmin)
	var err error

	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}

	//get result from insert
	user, err := result.Result()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepositories) FindUserByEmail(email string) (*users.UserCredentialsCheck, error) {

	query := `
		SELECT "id", "email", "password", "username", "role_id"
		FROM "users"
		WHERE "email" = $1
		`

	user := new(users.UserCredentialsCheck)
	if err := r.db.Get(user, query, email); err != nil {
		return nil, fmt.Errorf("email not found")
	}

	return user, nil

}

func (r *userRepositories) InsertOauth(req *users.UserPassport) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	INSERT INTO "oauth" (
		"user_id",
		"refresh_tokens",
		"access_tokens"
	)
	VALUES ($1, $2, $3)
		RETURNING "id"
	`
	if err := r.db.QueryRowContext(
		ctx,
		query,
		req.User.ID,
		req.Token.RefreshToken,
		req.Token.AccessToken,
	).Scan(&req.Token.Id); err != nil {
		return fmt.Errorf("insert oauth failed: %v", err)
	}
	return nil
}

func (r *userRepositories) FindOneOauth(refreshToken string) (*users.OAuth, error) {
	query := `
		SELECT 
			"id",
			"user_id"
		FROM "oauth"
		WHERE "refresh_tokens" = $1;
	`
	oauth := new(users.OAuth)
	if err := r.db.Get(oauth, query, refreshToken); err != nil {
		return nil, fmt.Errorf("oauth not found %v", err)
	}
	return oauth, nil

}

func (r *userRepositories) UpdateOauth(req *users.UserTokens) error {
	query := `
		UPDATE "oauth" SET
			"access_tokens" = :access_tokens,
			"refresh_tokens" = :refresh_tokens
		WHERE "id" = :id;
		`

	if _, err := r.db.NamedExecContext(context.Background(), query, req); err != nil {
		return fmt.Errorf("update oauth failed: %v", err)
	}
	return nil
}

func (r *userRepositories) GetProfile(userId string) (*users.User, error) {

	query := `
		SELECT
			"id",
			"email",
			"username",
			"role_id"
		FROM "users"
		WHERE "id" = $1;
	`

	profile := new(users.User)
	if err := r.db.Get(profile, query, userId); err != nil {
		return nil, fmt.Errorf("profile not found: %v", err)
	}
	return profile, nil

}

func (r *userRepositories) DeleteOAuth(oauthId string) error {
	query := `
		DELETE FROM "oauth" WHERE "id" = $1;
	`

	if _, err := r.db.ExecContext(context.Background(), query, oauthId); err != nil {
		return fmt.Errorf("oauth not found")
	}

	return nil
}

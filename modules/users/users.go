package users

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Username string `db:"username" json:"username"`
	RoleId   int    `db:"role_id" json:"role_id"`
}

type UserRegisterRequest struct {
	Email    string `db:"email" json:"email" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
	Username string `db:"username" json:"username" form:"username"`
}

func (obj *UserRegisterRequest) BcryptHashing() error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 10)
	if err != nil {
		return fmt.Errorf("bcrypt hashing error: %v", err)
	}
	obj.Password = string(hashPassword)
	return nil
}

func (obj *UserRegisterRequest) ValidateEmail() bool {
	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, obj.Email)
	if err != nil {
		return false
	}
	return match
}

type UserPassport struct {
	User  *User       `json:"user"`
	Token *UserTokens `json:"token"`
}

type UserTokens struct {
	Id           string `db:"id" json:"id"`
	AccessToken  string `db:"access_tokens" json:"access_tokens"`
	RefreshToken string `db:"refresh_tokens" json:"refresh_tokens"`
}

type UserCredentials struct {
	Email    string `db:"email" json:"email" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
}

type UserCredentialsCheck struct {
	Id       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Username string `db:"username" json:"username"`
	RoleId   int    `db:"role_id" json:"role_id"`
}

type UserRefreshCredentials struct {
	RefreshToken string `db:"refresh_tokens" json:"refresh_tokens" form:"refresh_tokens"`
}

type UserClaims struct {
	Id     string `db:"id" json:"id"`
	RoleId int    `db:"role_id" json:"role_id"`
}

type OAuth struct {
	Id     string `db:"id" json:"id"`
	UserId string `db:"user_id" json:"user_id"`
}

type UserRemoveCredentials struct {
	OauthId string `db:"id" json:"id" form:"id"`
}

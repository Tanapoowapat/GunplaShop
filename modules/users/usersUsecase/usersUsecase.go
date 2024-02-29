package usersUsecase

import (
	"fmt"

	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/Tanapoowapat/GunplaShop/modules/users"
	"github.com/Tanapoowapat/GunplaShop/modules/users/usersRepositories"
	"github.com/Tanapoowapat/GunplaShop/pkg/gunplaauth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterRequest) (*users.UserPassport, error)
	InsertAdmin(req *users.UserRegisterRequest) (*users.UserPassport, error)
	GetUserProfile(userId string) (*users.User, error)
	GetPassport(req *users.UserCredentials) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredentials) (*users.UserPassport, error)
	DeleteOAuth(oauthId string) error
}

type usersUsecase struct {
	config    config.IConfig
	user_repo usersRepositories.IUserRepositories
}

func UsersUsecase(config config.IConfig, user_repo usersRepositories.IUserRepositories) IUsersUsecase {
	return &usersUsecase{
		config:    config,
		user_repo: user_repo,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterRequest) (*users.UserPassport, error) {
	// Hashing password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}
	// Insert User
	result, err := u.user_repo.InsertUser(req, false)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (u *usersUsecase) InsertAdmin(req *users.UserRegisterRequest) (*users.UserPassport, error) {
	// Hashing password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}
	// Insert User
	result, err := u.user_repo.InsertUser(req, true)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *usersUsecase) GetPassport(req *users.UserCredentials) (*users.UserPassport, error) {
	// Find User
	user, err := u.user_repo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Signin JWT
	access_token, err := gunplaauth.NewAuthTokens(gunplaauth.AccessToken, u.config.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, err
	}

	//Refresh JWT
	refresh_token, err := gunplaauth.NewAuthTokens(gunplaauth.RefeshToken, u.config.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, err
	}
	// return user Passport
	passort := &users.UserPassport{
		User: &users.User{
			ID:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserTokens{
			AccessToken:  access_token.SignToken(),
			RefreshToken: refresh_token.SignToken(),
		},
	}

	if err := u.user_repo.InsertOauth(passort); err != nil {
		return nil, err
	}

	return passort, nil
}

func (u *usersUsecase) RefreshPassport(req *users.UserRefreshCredentials) (*users.UserPassport, error) {
	//Parse token
	claims, err := gunplaauth.ParseToken(u.config.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	//Find Oauth
	oauth, err := u.user_repo.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	profile, err := u.user_repo.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:     profile.ID,
		RoleId: profile.RoleId,
	}

	accessToken, err := gunplaauth.NewAuthTokens(gunplaauth.AccessToken, u.config.Jwt(), newClaims)
	if err != nil {
		return nil, err
	}

	refresh_token := gunplaauth.RepeatyToken(
		u.config.Jwt(),
		newClaims,
		claims.ExpiresAt.Unix(),
	)

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserTokens{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refresh_token,
		},
	}
	if err := u.user_repo.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u *usersUsecase) GetUserProfile(userId string) (*users.User, error) {
	profile, err := u.user_repo.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (u *usersUsecase) DeleteOAuth(oauthId string) error {
	if err := u.user_repo.DeleteOAuth(oauthId); err != nil {
		return err
	}
	return nil
}

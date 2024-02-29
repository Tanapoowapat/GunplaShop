package middlewaresUsecase

import (
	"github.com/Tanapoowapat/GunplaShop/modules/middlewares"
	"github.com/Tanapoowapat/GunplaShop/modules/middlewares/middlewaresRepositories"
)

type IMiddlewaresUsecase interface {
	FindAcessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresUsecase struct {
	repo middlewaresRepositories.IMiddlewaresRepositories
}

func NewMiddlewaresUsecase(repo middlewaresRepositories.IMiddlewaresRepositories) IMiddlewaresUsecase {
	return &middlewaresUsecase{
		repo: repo,
	}
}

func (mu *middlewaresUsecase) FindAcessToken(userId, accessToken string) bool {
	return mu.repo.FindAcessToken(userId, accessToken)
}

func (mu *middlewaresUsecase) FindRole() ([]*middlewares.Role, error) {
	roles, err := mu.repo.FindRole()
	if err != nil {
		return nil, err
	}
	return roles, nil

}

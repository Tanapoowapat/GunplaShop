package middlewaresUsecase

import "github.com/Tanapoowapat/GunplaShop/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecase interface{}

type middlewaresUsecase struct {
	repo middlewaresRepositories.IMiddlewaresRepositories
}

func NewMiddlewaresUsecase(repo middlewaresRepositories.IMiddlewaresRepositories) IMiddlewaresUsecase {
	return &middlewaresUsecase{
		repo: repo,
	}
}

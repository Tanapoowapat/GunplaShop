package middlewaresHandlers

import (
	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/Tanapoowapat/GunplaShop/modules/entities"
	"github.com/Tanapoowapat/GunplaShop/modules/middlewares/middlewaresUsecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewaresHandlersErrorCode string

const (
	routerCheckErr middlewaresHandlersErrorCode = "middlewares-001"
)

type IMiddlewaresHandlers interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
}

type middlewaresHandlers struct {
	config  config.IConfig
	usecase middlewaresUsecase.IMiddlewaresUsecase
}

func NewMiddlewaresHandlers(cfg config.IConfig, usecase middlewaresUsecase.IMiddlewaresUsecase) IMiddlewaresHandlers {
	return &middlewaresHandlers{
		config:  cfg,
		usecase: usecase,
	}
}

func (mh *middlewaresHandlers) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (mh *middlewaresHandlers) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(fiber.ErrNotFound.Code, string(routerCheckErr), "Router not found").Res()
	}
}

func (mh *middlewaresHandlers) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Asia/Bangkok",
	})
}

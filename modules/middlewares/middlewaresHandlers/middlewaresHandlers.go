package middlewaresHandlers

import (
	"strings"

	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/Tanapoowapat/GunplaShop/modules/entities"
	"github.com/Tanapoowapat/GunplaShop/modules/middlewares/middlewaresUsecase"
	"github.com/Tanapoowapat/GunplaShop/pkg/gunplaauth"
	"github.com/Tanapoowapat/GunplaShop/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewaresHandlersErrorCode string

const (
	routerCheckErr   middlewaresHandlersErrorCode = "middlewares-001"
	JwtAuthErr       middlewaresHandlersErrorCode = "middlewares-002"
	ParamsCheckErr   middlewaresHandlersErrorCode = "middlewares-003"
	authorizationErr middlewaresHandlersErrorCode = "middlewares-004"
)

type IMiddlewaresHandlers interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorization(expectRoleId ...int) fiber.Handler
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

func (mh *middlewaresHandlers) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := gunplaauth.ParseToken(mh.config.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.StatusUnauthorized,
				string(JwtAuthErr),
				"Unauthorized").Res()
		}
		claims := result.Claims
		if !mh.usecase.FindAcessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.StatusUnauthorized,
				string(JwtAuthErr),
				"No Permission to access").Res()
		}
		// Set user id to locals
		c.Locals("userId", claims.Id)
		c.Locals("userRoleID", claims.RoleId)
		return c.Next()
	}
}

func (mh *middlewaresHandlers) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Params("userId") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(ParamsCheckErr),
				"Unauthorized",
			).Res()
		}
		return c.Next()
	}
}

func (mh *middlewaresHandlers) Authorization(expectRoleId ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals("userRoleID").(int)
		if !ok {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizationErr),
				"User id is not Int Type",
			).Res()
		}

		roles, err := mh.usecase.FindRole()
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(authorizationErr),
				err.Error(),
			).Res()
		}

		sum := 0
		for _, v := range expectRoleId {
			sum += v
		}

		expectValbin := utils.BinaryConverter(sum, len(roles))
		userValbin := utils.BinaryConverter(userRoleId, len(roles))

		for i := range userValbin {
			if userValbin[i]&expectValbin[i] == 1 {
				return c.Next()
			}
		}

		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizationErr),
			"Unauthorized",
		).Res()
	}
}

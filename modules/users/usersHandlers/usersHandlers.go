package usersHandlers

import (
	"strings"

	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/Tanapoowapat/GunplaShop/modules/entities"
	"github.com/Tanapoowapat/GunplaShop/modules/users"
	"github.com/Tanapoowapat/GunplaShop/modules/users/usersUsecase"
	"github.com/Tanapoowapat/GunplaShop/pkg/gunplaauth"
	"github.com/gofiber/fiber/v2"
)

type usersHandlersErrCode string

const (
	signUpCustomerErrCode     usersHandlersErrCode = "user-001-1"
	signUpAdminErrCode        usersHandlersErrCode = "user-001-2"
	signInErrCode             usersHandlersErrCode = "user-002"
	refreshPassportErrCode    usersHandlersErrCode = "user-003"
	signOutErrCode            usersHandlersErrCode = "user-004"
	generateAdminTokenErrCode usersHandlersErrCode = "user-005"
	getUserprofileErrCode     usersHandlersErrCode = "user-006"
)

type IUserHandlers interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignUpAdmin(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
	GenaerateAdminToken(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
}

type usersHandlers struct {
	cfg         config.IConfig
	userUsecase usersUsecase.IUsersUsecase
}

func NewUsersHandlers(cfg config.IConfig, userUsecase usersUsecase.IUsersUsecase) IUserHandlers {
	return &usersHandlers{
		cfg:         cfg,
		userUsecase: userUsecase,
	}
}

func (h *usersHandlers) SignUpCustomer(c *fiber.Ctx) error {
	req := new(users.UserRegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signUpCustomerErrCode),
			err.Error(),
		).Res()
	}

	// Validate Email
	if !req.ValidateEmail() {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signUpCustomerErrCode),
			"Invalid Email",
		).Res()
	}

	//Insert User
	result, err := h.userUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(signUpCustomerErrCode),
				"Username has been used",
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(signUpCustomerErrCode),
				"Email has been used",
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.StatusInternalServerError,
				string(signUpCustomerErrCode),
				err.Error(),
			).Res()

		}
	}

	return entities.NewResponse(c).Sucess(fiber.StatusCreated, result).Res()
}

func (h *usersHandlers) SignUpAdmin(c *fiber.Ctx) error {
	req := new(users.UserRegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signUpCustomerErrCode),
			err.Error(),
		).Res()
	}

	// Validate Email
	if !req.ValidateEmail() {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signUpCustomerErrCode),
			"Invalid Email",
		).Res()
	}
	//Insert User
	result, err := h.userUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(signUpCustomerErrCode),
				"Username has been used",
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(signUpCustomerErrCode),
				"Email has been used",
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.StatusInternalServerError,
				string(signUpCustomerErrCode),
				err.Error(),
			).Res()
		}
	}

	return entities.NewResponse(c).Sucess(fiber.StatusCreated, result).Res()
}

func (h *usersHandlers) SignIn(c *fiber.Ctx) error {

	// Parse Body
	user := new(users.UserCredentials)
	if err := c.BodyParser(user); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signInErrCode),
			err.Error(),
		).Res()
	}

	// Get Passport
	passort, err := h.userUsecase.GetPassport(user)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signInErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Sucess(fiber.StatusOK, passort).Res()
}

func (h *usersHandlers) RefreshPassport(c *fiber.Ctx) error {

	// Parse Body
	req := new(users.UserRefreshCredentials)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(refreshPassportErrCode),
			err.Error(),
		).Res()
	}

	// //Get Passport
	passort, err := h.userUsecase.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(refreshPassportErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Sucess(fiber.StatusOK, passort).Res()

}

func (h *usersHandlers) SignOut(c *fiber.Ctx) error {

	// Parse Body
	req := new(users.UserRemoveCredentials)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signOutErrCode),
			err.Error(),
		).Res()
	}

	if err := h.userUsecase.DeleteOAuth(req.OauthId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(signOutErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Sucess(fiber.StatusOK, nil).Res()

}

func (h *usersHandlers) GenaerateAdminToken(c *fiber.Ctx) error {
	adminToken, err := gunplaauth.NewAuthTokens(gunplaauth.AdminToken, h.cfg.Jwt(), nil)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateAdminTokenErrCode),
			err.Error()).Res()
	}

	return entities.NewResponse(c).Sucess(
		fiber.StatusOK, &struct {
			Token string `json:"token"`
		}{
			Token: adminToken.SignToken(),
		}).Res()
}

func (h *usersHandlers) GetUserProfile(c *fiber.Ctx) error {
	//set params
	userId := strings.Trim(c.Params("userId"), " ")

	result, err := h.userUsecase.GetUserProfile(userId)
	if err != nil {
		switch err.Error() {
		case "profile not found: sql: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(getUserprofileErrCode),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.StatusInternalServerError,
				string(getUserprofileErrCode),
				err.Error(),
			).Res()
		}
	}

	return entities.NewResponse(c).Sucess(fiber.StatusOK, result).Res()
}

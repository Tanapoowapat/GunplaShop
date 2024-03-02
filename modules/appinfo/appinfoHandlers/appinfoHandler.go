package appinfohandlers

import (
	"strconv"
	"strings"

	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/Tanapoowapat/GunplaShop/modules/appinfo"
	appinfousecase "github.com/Tanapoowapat/GunplaShop/modules/appinfo/appinfoUsecase"
	"github.com/Tanapoowapat/GunplaShop/modules/entities"
	"github.com/Tanapoowapat/GunplaShop/pkg/gunplaauth"
	"github.com/gofiber/fiber/v2"
)

type appinfoHandlersErrCode string

const (
	generateApiKeyErrCode appinfoHandlersErrCode = "appinfo-001"
	findCategoryErrCode   appinfoHandlersErrCode = "appinfo-002"
	InsertCategoryErrCode appinfoHandlersErrCode = "appinfo-003"
	DeleteCategoryErrCode appinfoHandlersErrCode = "appinfo-004"
)

type IAppinfoHanlder interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
	InsertCategory(c *fiber.Ctx) error
	RemoveCategory(c *fiber.Ctx) error
}

type appinfoHandlers struct {
	cfg             config.IConfig
	appinfo_usecase appinfousecase.IAppinfoUsecase
}

func NewAppinfoHandlers(cfg config.IConfig, appinfo_usecase appinfousecase.IAppinfoUsecase) IAppinfoHanlder {
	return &appinfoHandlers{
		cfg:             cfg,
		appinfo_usecase: appinfo_usecase,
	}
}

func (h *appinfoHandlers) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := gunplaauth.NewAuthTokens(
		gunplaauth.ApiKeyToken,
		h.cfg.Jwt(),
		nil,
	)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(generateApiKeyErrCode),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Sucess(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		}{
			Key: apiKey.SignToken(),
		},
	).Res()

}

func (h *appinfoHandlers) FindCategory(c *fiber.Ctx) error {
	req := new(appinfo.CategoryFiter)
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryErrCode),
			err.Error(),
		).Res()
	}

	cat, err := h.appinfo_usecase.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryErrCode),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Sucess(fiber.StatusOK, cat).Res()

}

func (h *appinfoHandlers) InsertCategory(c *fiber.Ctx) error {
	req := make([]*appinfo.Category, 0)

	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(InsertCategoryErrCode),
			err.Error(),
		).Res()
	}
	if len(req) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(InsertCategoryErrCode),
			"Category request cannot be empty",
		).Res()
	}

	if err := h.appinfo_usecase.InsertCategory(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(InsertCategoryErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Sucess(fiber.StatusCreated, req).Res()
}

func (h *appinfoHandlers) RemoveCategory(c *fiber.Ctx) error {
	category_id := strings.Trim(c.Params("category_id"), " ")
	categoryId, err := strconv.Atoi(category_id)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(DeleteCategoryErrCode),
			err.Error(),
		).Res()
	}

	if categoryId <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(DeleteCategoryErrCode),
			"Invalid Category Id",
		).Res()
	}

	if err := h.appinfo_usecase.DeleteCategory(categoryId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(DeleteCategoryErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Sucess(
		fiber.StatusCreated,
		&struct {
			Category_id string `json:"category_id"`
		}{
			Category_id: category_id,
		}).Res()
}

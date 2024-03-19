package filehandler

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/Tanapoowapat/GunplaShop/modules/entities"
	"github.com/Tanapoowapat/GunplaShop/modules/file"
	filesusecase "github.com/Tanapoowapat/GunplaShop/modules/file/filesUsecase"
	"github.com/Tanapoowapat/GunplaShop/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type fileHandlerErr string

const (
	UploadErr fileHandlerErr = "File-001"
	DeleteErr fileHandlerErr = "File-002"
)

type IFileHandler interface {
	UploadFile(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
}

type fileHandler struct {
	cfg     config.IConfig
	usecase filesusecase.IFileUsecase
}

func NewfileHandler(cfg config.IConfig, fileusecase filesusecase.IFileUsecase) IFileHandler {
	return &fileHandler{
		cfg:     cfg,
		usecase: fileusecase,
	}
}

func (h *fileHandler) UploadFile(c *fiber.Ctx) error {

	req := make([]*file.FileReq, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(UploadErr),
			err.Error(),
		).Res()
	}

	filesReq := form.File["file"]
	destination := c.FormValue("destination")

	// File Validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}

	for _, f := range filesReq {
		ext := strings.TrimPrefix(filepath.Ext(f.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(UploadErr),
				"Invalid File",
			).Res()
		}

		if f.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(UploadErr),
				fmt.Sprintf("File must be less than %d MB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}
		filename := utils.RandomFileName(ext)
		req = append(req, &file.FileReq{
			File:        f,
			Destination: destination + "/" + filename,
			FileName:    filename,
			Extension:   ext,
		})
	}

	res, err := h.usecase.UploadImageLocal(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(UploadErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Sucess(fiber.StatusOK, res).Res()
}

func (h *fileHandler) DeleteFile(c *fiber.Ctx) error {

	req := make([]*file.DeleteFileReq, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(DeleteErr),
			err.Error(),
		).Res()
	}

	if err := h.usecase.DeleteImageStorage(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(DeleteErr),
			err.Error(),
		).Res()
	}

	return nil
}

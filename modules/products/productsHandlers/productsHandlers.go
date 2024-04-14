package productshandlers

import (
	"fmt"
	"strings"

	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/Tanapoowapat/GunplaShop/modules/appinfo"
	"github.com/Tanapoowapat/GunplaShop/modules/entities"
	"github.com/Tanapoowapat/GunplaShop/modules/file"
	filesusecase "github.com/Tanapoowapat/GunplaShop/modules/file/filesUsecase"
	"github.com/Tanapoowapat/GunplaShop/modules/products"
	productsusecase "github.com/Tanapoowapat/GunplaShop/modules/products/productsUsercase"
	"github.com/gofiber/fiber/v2"
)

type productsHandlerErr string

const (
	FindOneErr       productsHandlerErr = "Products-001"
	FindProductErr   productsHandlerErr = "Products-002"
	AddProductErr    productsHandlerErr = "Products-003"
	DeleteProductErr productsHandlerErr = "Products-004"
	UpdateProductErr productsHandlerErr = "Products-005"
)

type IProductsHandler interface {
	FindOneProduct(c *fiber.Ctx) error
	FindProducts(c *fiber.Ctx) error
	AddProducts(c *fiber.Ctx) error
	DeleteProducts(c *fiber.Ctx) error
	UpdateProducts(c *fiber.Ctx) error
}

type productsHandler struct {
	cfg         config.IConfig
	prodUsecase productsusecase.IProductUseCase
	fileUsecase filesusecase.IFileUsecase
}

func NewProductsHandler(cfg config.IConfig, prodUsecase productsusecase.IProductUseCase, fileUsecase filesusecase.IFileUsecase) IProductsHandler {
	return &productsHandler{
		cfg:         cfg,
		prodUsecase: prodUsecase,
		fileUsecase: fileUsecase,
	}
}

func (handler *productsHandler) FindOneProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")
	product, err := handler.prodUsecase.FindOneProducts(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(FindOneErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Sucess(
		fiber.StatusOK,
		product,
	).Res()
}

func (handler *productsHandler) FindProducts(c *fiber.Ctx) error {

	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(FindOneErr),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 5 {
		req.Limit = 5
	}

	if req.OrderBy == "" {
		req.OrderBy = "title"
	}

	if req.Sort == "" {
		req.Sort = "ASC"
	}

	products := handler.prodUsecase.FindProducts(req)

	return entities.NewResponse(c).Sucess(
		fiber.StatusOK,
		products,
	).Res()
}

func (h *productsHandler) AddProducts(c *fiber.Ctx) error {
	req := &products.Products{
		Category: &appinfo.Category{},
		Images:   []*entities.Images{},
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddProductErr),
			err.Error(),
		).Res()
	}

	if req.Category.Id <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(FindOneErr),
			"Category id Invalid",
		).Res()
	}

	product, err := h.prodUsecase.AddProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(AddProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Sucess(fiber.StatusOK, product).Res()

}

func (h *productsHandler) DeleteProducts(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.prodUsecase.FindOneProducts(productId)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(DeleteProductErr), err.Error()).Res()
	}

	deleteFileReq := make([]*file.DeleteFileReq, 0)
	for _, p := range product.Images {
		deleteFileReq = append(deleteFileReq, &file.DeleteFileReq{
			Destination: fmt.Sprintf("images/products/%s", p.FileName),
		})
	}

	//Delete From Local
	if err := h.fileUsecase.DeleteImageLocal(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(DeleteProductErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Sucess(fiber.StatusNoContent, nil).Res()

}

func (h *productsHandler) UpdateProducts(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")
	req := &products.Products{
		Images:   make([]*entities.Images, 0),
		Category: &appinfo.Category{},
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadGateway.Code,
			string(AddProductErr),
			err.Error(),
		).Res()
	}
	req.Id = productId

	product, err := h.prodUsecase.UpdateProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(AddProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Sucess(fiber.StatusOK, product).Res()

}

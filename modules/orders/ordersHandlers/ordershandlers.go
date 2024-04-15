package ordershandlers

import (
	"strings"
	"time"

	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/Tanapoowapat/GunplaShop/modules/entities"
	"github.com/Tanapoowapat/GunplaShop/modules/orders"
	ordersusecase "github.com/Tanapoowapat/GunplaShop/modules/orders/ordersUsecase"
	"github.com/gofiber/fiber/v2"
)

type OrdersHandlersErr string

const (
	FindOnceOrdersErr OrdersHandlersErr = "Orders-001"
	FindOrdersErr     OrdersHandlersErr = "Orders-002"
)

type IOrdersHandlers interface {
	FindOnceOrders(c *fiber.Ctx) error
	FindOrder(c *fiber.Ctx) error
}

type ordersHandlers struct {
	ordersUsecase ordersusecase.IOrdersUsecase
	cfg           config.IConfig
}

func NewOrdersHandlers(ordersUsecase ordersusecase.IOrdersUsecase, cfg config.IConfig) IOrdersHandlers {
	return &ordersHandlers{
		ordersUsecase: ordersUsecase,
		cfg:           cfg,
	}
}

func (handler *ordersHandlers) FindOnceOrders(c *fiber.Ctx) error {

	orderId := strings.Trim(c.Params("order_id"), " ")
	order, err := handler.ordersUsecase.FindOnceOrders(orderId)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code,
			string(FindOnceOrdersErr),
			err.Error()).Res()
	}

	return entities.NewResponse(c).Sucess(fiber.StatusOK, order).Res()
}

func (handler *ordersHandlers) FindOrder(c *fiber.Ctx) error {
	req := &orders.OrderFilter{
		SortReq:       &entities.SortReq{},
		PaginationReq: &entities.PaginationReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(FindOrdersErr),
			err.Error(),
		).Res()
	}

	// Paginate
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 5
	}

	// Sort
	orderByMap := map[string]string{
		"id":         `"o"."id"`,
		"created_at": `"o"."created_at"`,
	}
	if orderByMap[req.OrderBy] == "" {
		req.OrderBy = orderByMap["id"]
	}

	req.Sort = strings.ToUpper(req.Sort)
	sortMap := map[string]string{
		"DESC": "DESC",
		"ASC":  "ASC",
	}
	if sortMap[req.Sort] == "" {
		req.Sort = sortMap["DESC"]
	}

	// Date	YYYY-MM-DD
	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(FindOrdersErr),
				"start date is invalid",
			).Res()
		}
		req.StartDate = start.Format("2006-01-02")
	}

	// Date	YYYY-MM-DD
	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(FindOrdersErr),
				"end date is invalid",
			).Res()
		}
		req.EndDate = end.Format("2006-01-02")
	}

	return entities.NewResponse(c).Sucess(
		fiber.StatusOK,
		handler.ordersUsecase.FindOrders(req),
	).Res()
}

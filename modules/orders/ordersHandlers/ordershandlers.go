package ordershandlers

import (
	"strings"
	"time"

	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/Tanapoowapat/GunplaShop/modules/entities"
	"github.com/Tanapoowapat/GunplaShop/modules/orders"
	ordersusecase "github.com/Tanapoowapat/GunplaShop/modules/orders/ordersUsecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrdersHandlersErr string

const (
	FindOnceOrdersErr OrdersHandlersErr = "Orders-001"
	FindOrdersErr     OrdersHandlersErr = "Orders-002"
	InsertOrderErr    OrdersHandlersErr = "Orders-003"
	UpdateOrderErr    OrdersHandlersErr = "Orders-004"
)

type IOrdersHandlers interface {
	FindOnceOrders(c *fiber.Ctx) error
	FindOrder(c *fiber.Ctx) error
	InsertOrder(c *fiber.Ctx) error
	UpdateOrder(c *fiber.Ctx) error
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

func (h *ordersHandlers) InsertOrder(c *fiber.Ctx) error {
	userId := c.Locals("usersId").(string)

	req := &orders.Order{
		Product: make([]*orders.ProductOrder, 0),
	}

	if err := c.BodyParser(req); err != nil {
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(InsertOrderErr),
				err.Error(),
			).Res()
		}
	}

	if len(req.Product) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(InsertOrderErr),
			"product are empty",
		).Res()
	}

	if c.Locals("userRoleId").(int) != 2 {
		req.UserId = userId
	}

	req.Status = "waiting"
	req.TotalPrice = 0

	order, err := h.ordersUsecase.InsertOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(FindOrdersErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Sucess(fiber.StatusOK, order).Res()
}

func (h *ordersHandlers) UpdateOrder(c *fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")
	req := new(orders.Order)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(UpdateOrderErr), err.Error()).Res()
	}

	req.Id = orderId
	statusMap := map[string]string{
		"waiting":   "waiting",
		"shipping":  "shipping",
		"completed": "completed",
		"canceled":  "canceled",
	}

	//Check User in Local Cache
	if c.Locals("userRoleId").(int) == 2 {
		req.Status = statusMap[strings.ToLower(req.Status)]
	} else if strings.ToLower(req.Status) == statusMap["canceled"] {
		req.Status = statusMap["canceled"]
	}

	if req.TransferSlip != nil {
		if req.TransferSlip.Id == "" {
			req.TransferSlip.Id = uuid.NewString()
		}
		if req.TransferSlip.CreatedAt == "" {
			loc, err := time.LoadLocation("Asia/Bangkok")
			if err != nil {
				return entities.NewResponse(c).Error(
					fiber.ErrInternalServerError.Code,
					string(UpdateOrderErr),
					err.Error(),
				).Res()
			}
			now := time.Now().In(loc)
			// YYYY-MM-DD HH:MM:SS
			// 2006-01-02 15:04:05
			req.TransferSlip.CreatedAt = now.Format("2006-01-02 15:04:05")

		}
	}

	order, err := h.ordersUsecase.UpdateOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(UpdateOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Sucess(fiber.StatusOK, order).Res()
}

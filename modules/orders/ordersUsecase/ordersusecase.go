package ordersusecase

import (
	"math"

	"github.com/Tanapoowapat/GunplaShop/modules/entities"
	"github.com/Tanapoowapat/GunplaShop/modules/orders"
	ordersrepositories "github.com/Tanapoowapat/GunplaShop/modules/orders/ordersRepositories"
	productsrepositories "github.com/Tanapoowapat/GunplaShop/modules/products/productsRepositories"
)

type IOrdersUsecase interface {
	FindOnceOrders(orderId string) (*orders.Order, error)
	FindOrders(req *orders.OrderFilter) *entities.PaginateRes
}

type ordersUsecase struct {
	ordersRepo   ordersrepositories.IOrdersRepositories
	productsRepo productsrepositories.IProductRepositorise
}

func NewOrdersUsecase(ordersRepo ordersrepositories.IOrdersRepositories, productsRepo productsrepositories.IProductRepositorise) IOrdersUsecase {
	return &ordersUsecase{
		ordersRepo:   ordersRepo,
		productsRepo: productsRepo,
	}
}

func (usecase *ordersUsecase) FindOnceOrders(orderId string) (*orders.Order, error) {
	return usecase.ordersRepo.FindOnceOrders(orderId)
}

func (usecase *ordersUsecase) FindOrders(req *orders.OrderFilter) *entities.PaginateRes {
	order, count := usecase.ordersRepo.FindOrders(req)
	return &entities.PaginateRes{
		Data:       order,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalItems: count,
		TotalPage:  int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

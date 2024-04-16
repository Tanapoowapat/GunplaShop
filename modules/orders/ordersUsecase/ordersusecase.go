package ordersusecase

import (
	"fmt"
	"math"

	"github.com/Tanapoowapat/GunplaShop/modules/entities"
	"github.com/Tanapoowapat/GunplaShop/modules/orders"
	ordersrepositories "github.com/Tanapoowapat/GunplaShop/modules/orders/ordersRepositories"
	productsrepositories "github.com/Tanapoowapat/GunplaShop/modules/products/productsRepositories"
)

type IOrdersUsecase interface {
	FindOnceOrders(orderId string) (*orders.Order, error)
	FindOrders(req *orders.OrderFilter) *entities.PaginateRes
	InsertOrder(req *orders.Order) (*orders.Order, error)
	UpdateOrder(req *orders.Order) (*orders.Order, error)
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

func (usecase *ordersUsecase) InsertOrder(req *orders.Order) (*orders.Order, error) {
	//check if product exits
	for i := range req.Product {
		if req.Product[i].Product == nil {
			return nil, fmt.Errorf("product is empty")
		}

		//Find Product
		product, err := usecase.productsRepo.FindOneProducts(req.Product[i].Product.Id)
		if err != nil {
			return nil, err
		}

		//Set price
		req.TotalPrice += req.Product[i].Product.Price * float64(req.Product[i].Qty)
		req.Product[i].Product = product
	}

	orderId, err := usecase.ordersRepo.InsertOrder(req)
	if err != nil {
		return nil, err
	}

	order, err := usecase.FindOnceOrders(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil

}

func (u *ordersUsecase) UpdateOrder(req *orders.Order) (*orders.Order, error) {

	if err := u.ordersRepo.UpdateOrder(req); err != nil {
		return nil, err
	}

	order, err := u.ordersRepo.FindOnceOrders(req.Id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

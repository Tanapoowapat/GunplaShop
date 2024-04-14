package productsusecase

import (
	"math"

	"github.com/Tanapoowapat/GunplaShop/modules/entities"
	"github.com/Tanapoowapat/GunplaShop/modules/products"
	productsrepositories "github.com/Tanapoowapat/GunplaShop/modules/products/productsRepositories"
)

type IProductUseCase interface {
	FindOneProducts(productId string) (*products.Products, error)
	FindProducts(req *products.ProductFilter) *entities.PaginateRes
	AddProduct(req *products.Products) (*products.Products, error)
	DeleteProduct(productId string) error
	UpdateProduct(req *products.Products) (*products.Products, error)
}

type productsUsecase struct {
	productsRepo productsrepositories.IProductRepositorise
}

func NewProductsUsecase(productsRepo productsrepositories.IProductRepositorise) IProductUseCase {
	return &productsUsecase{
		productsRepo: productsRepo,
	}
}

func (usecase *productsUsecase) FindOneProducts(productId string) (*products.Products, error) {
	product, err := usecase.productsRepo.FindOneProducts(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (usecase *productsUsecase) FindProducts(req *products.ProductFilter) *entities.PaginateRes {

	data, count := usecase.productsRepo.FindProduct(req)

	return &entities.PaginateRes{
		Data:       data,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalItems: count,
		TotalPage:  int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (usecase *productsUsecase) AddProduct(req *products.Products) (*products.Products, error) {
	product, err := usecase.productsRepo.InsertProduct(req)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (usecase *productsUsecase) DeleteProduct(productId string) error {
	return usecase.productsRepo.DeleteProduct(productId)
}

func (usecase *productsUsecase) UpdateProduct(req *products.Products) (*products.Products, error) {
	return usecase.productsRepo.UpdateProduct(req)
}

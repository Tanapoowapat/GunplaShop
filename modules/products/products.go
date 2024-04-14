package products

import (
	"github.com/Tanapoowapat/GunplaShop/modules/appinfo"
	"github.com/Tanapoowapat/GunplaShop/modules/entities"
)

type Products struct {
	Id          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Category    *appinfo.Category  `json:"category"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
	Price       float64            `json:"price"`
	Images      []*entities.Images `json:"media"`
}

type ProductFilter struct {
	Id     string `query:"id"`
	Search string `query:"search"` // title & description
	*entities.PaginationReq
	*entities.SortReq
}

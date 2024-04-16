package orders

import (
	"github.com/Tanapoowapat/GunplaShop/modules/entities"
	"github.com/Tanapoowapat/GunplaShop/modules/products"
)

type OrderFilter struct {
	Search    string `query:"search"`
	Status    string `query:"status"`
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	*entities.PaginationReq
	*entities.SortReq
}

type Order struct {
	Id           string          `db:"id" json:"id"`
	UserId       string          `db:"user_id" json:"user_id"`
	TransferSlip *TransferSlip   `db:"transfer_slip" json:"transfer_slip"`
	Product      []*ProductOrder `json:"products"`
	Address      string          `db:"address" json:"address"`
	Contact      string          `db:"contact" json:"contact"`
	Status       string          `db:"status" json:"status"`
	TotalPrice   float64         `db:"total_price" json:"total_price"`
	CreatedAt    string          `db:"created_at" json:"created_at"`
	UpdatedAt    string          `db:"updated_at" json:"updated_at"`
}

type TransferSlip struct {
	Id        string `json:"id"`
	FileName  string `json:"filename"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

type ProductOrder struct {
	Id      string             `db:"id" json:"id"`
	Qty     int                `db:"qty" json:"qty"`
	Product *products.Products `db:"product" json:"product"`
}

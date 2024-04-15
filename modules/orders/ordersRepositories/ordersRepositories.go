package ordersrepositories

import (
	"encoding/json"
	"fmt"

	"github.com/Tanapoowapat/GunplaShop/modules/orders"
	"github.com/Tanapoowapat/GunplaShop/modules/orders/orderpattern"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepositories interface {
	FindOnceOrders(orderId string) (*orders.Order, error)
	FindOrders(req *orders.OrderFilter) ([]*orders.Order, int)
}

type ordersRepositories struct {
	db *sqlx.DB
}

func NewOrdersRepositories(db *sqlx.DB) IOrdersRepositories {
	return &ordersRepositories{
		db: db,
	}
}

func (repo *ordersRepositories) FindOnceOrders(orderId string) (*orders.Order, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"o"."id",
			"o"."user_id",
			"o"."transfer_slip",
			"o"."status",
			(
				SELECT
					array_to_json(array_agg("pt"))
				FROM (
					SELECT
						"spo"."id",
						"spo"."qty",
						"spo"."product"
					FROM "products_orders" "spo"
					WHERE "spo"."order_id" = "o"."id"
				) AS "pt"
			) AS "products",
			"o"."address",
			"o"."contact",
			(
				SELECT
					SUM(COALESCE(("po"."product"->>'price')::FLOAT*("po"."qty")::FLOAT, 0))
				FROM "products_orders" "po"
				WHERE "po"."order_id" = "o"."id"
			) AS "total_paid",
			"o"."created_at",
			"o"."updated_at"
		FROM "orders" "o"
		WHERE "o"."id" = $1
	) AS "t";`

	orderData := &orders.Order{
		TransferSlip: &orders.TransferSlip{},
		Product:      make([]*orders.ProductOrder, 0),
	}

	raw := make([]byte, 0)
	if err := repo.db.Get(&raw, query, orderId); err != nil {
		return nil, fmt.Errorf("error get order: %v", err)
	}

	if err := json.Unmarshal(raw, &orderData); err != nil {
		return nil, fmt.Errorf("error unmarshal json: %v", err)
	}

	return orderData, nil
}

func (repo *ordersRepositories) FindOrders(req *orders.OrderFilter) ([]*orders.Order, int) {
	builder := orderpattern.FindOrderBuilder(repo.db, req)
	engineer := orderpattern.FindOrderEngineer(builder)
	return engineer.FindOrder(), engineer.CountOrder()
}

package ordersrepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Tanapoowapat/GunplaShop/modules/orders"
	"github.com/Tanapoowapat/GunplaShop/modules/orders/orderpattern"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepositories interface {
	FindOnceOrders(orderId string) (*orders.Order, error)
	FindOrders(req *orders.OrderFilter) ([]*orders.Order, int)
	InsertOrder(req *orders.Order) (string, error)
	UpdateOrder(req *orders.Order) error
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

func (repo *ordersRepositories) InsertOrder(req *orders.Order) (string, error) {
	builder := orderpattern.NewInsertOrderBuilder(repo.db, req)
	orderId, err := orderpattern.NewInsertOrderEngineer(builder).InsertOrders()

	if err != nil {
		return "", err
	}
	return orderId, nil
}

func (repo *ordersRepositories) UpdateOrder(req *orders.Order) error {

	query := `UPDATE "orders" SET`

	queryWhereStack := make([]string, 0)
	values := make([]any, 0)
	lastIndex := 1

	if req.Status != "" {
		values = append(values, req.Status)
		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`"status" = $%d`, lastIndex))
		lastIndex++
	}

	if req.TransferSlip != nil {
		values = append(values, req.TransferSlip)
		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`"transfer_slip" = $%d`, lastIndex))
		lastIndex++
	}

	values = append(values, req.Id)
	queryClose := fmt.Sprintf(`WHERE "id" = $%d`, lastIndex)

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			query += strings.Replace(queryWhereStack[i], "?", ",", 1)
		} else {
			query += strings.Replace(queryWhereStack[i], "?", ",", 1)
		}
	}

	query += queryClose
	if _, err := repo.db.ExecContext(context.Background(), query, values...); err != nil {
		return fmt.Errorf("update order fail: %v", err)
	}

	return nil
}

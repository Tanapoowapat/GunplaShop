package orderpattern

import (
	"context"
	"fmt"
	"time"

	"github.com/Tanapoowapat/GunplaShop/modules/orders"
	"github.com/jmoiron/sqlx"
)

type IInsertOrderBuilder interface {
	initTransaction() error
	insertOrder() error
	insertProductOrder() error
	commit() error
	getOrdersId() string
}

type insertOrdersBuilder struct {
	db  *sqlx.DB
	req *orders.Order
	tx  *sqlx.Tx
}

func NewInsertOrderBuilder(db *sqlx.DB, req *orders.Order) IInsertOrderBuilder {
	return &insertOrdersBuilder{
		db:  db,
		req: req,
	}
}

func (b *insertOrdersBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *insertOrdersBuilder) insertOrder() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := `
	INSERT INTO "orders" (
		"user_id",
		"contact",
		"address",
		"transfer_slip",
		"status"
	)
	VALUES
	($1, $2, $3, $4, $5)
		RETURNING "id";`

	if err := b.tx.QueryRowxContext(ctx, query,
		b.req.UserId,
		b.req.Contact,
		b.req.Address,
		b.req.TransferSlip,
		b.req.Status,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert order fail: %v", err)
	}

	return nil

}

func (b *insertOrdersBuilder) insertProductOrder() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := `
	INSERT INTO "products_orders" (
		"order_id",
		"qty",
		"product"
	)
	VALUES`

	values := make([]any, 0)
	lastIndex := 0
	for i := range b.req.Product {
		values = append(values,
			b.req.Id,
			b.req.Product[i].Qty,
			b.req.Product[i].Product,
		)

		if i != len(b.req.Product)-1 {
			query += fmt.Sprintf(`
		($%d, $%d, $%d),`, lastIndex+1, lastIndex+2, lastIndex+3)
		} else {
			query += fmt.Sprintf(`
		($%d, $%d, $%d);`, lastIndex+1, lastIndex+2, lastIndex+3)
		}

		lastIndex += 3

	}

	if _, err := b.tx.ExecContext(ctx, query, values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("error insert order_product: %v", err)
	}

	return nil
}

func (b *insertOrdersBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (b *insertOrdersBuilder) getOrdersId() string { return b.req.Id }

type insertOrdersEngineer struct {
	builder IInsertOrderBuilder
}

func NewInsertOrderEngineer(b IInsertOrderBuilder) *insertOrdersEngineer {
	return &insertOrdersEngineer{
		builder: b,
	}
}

func (en *insertOrdersEngineer) InsertOrders() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", nil
	}
	if err := en.builder.insertOrder(); err != nil {
		return "", nil
	}
	if err := en.builder.insertProductOrder(); err != nil {
		return "", nil
	}
	if err := en.builder.commit(); err != nil {
		return "", nil
	}

	return en.builder.getOrdersId(), nil
}

package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
)

type OrderProduct struct {
	ID        int `db:"id" fieldtag:"pk" json:"id"`
	OrderID   int `db:"order_id" json:"orderId"`
	ProductID int `db:"product_id" json:"productId"`
}

type IOrderProductRepository interface {
	Create(ctx context.Context, tx pgx.Tx, orderID, productID int) error
}

type OrderProductRepository struct{}

const OrderProductsTable = "order_products"

var OrderProductStruct = sqlbuilder.NewStruct(new(OrderProduct))

func NewOrderProductRepository() *OrderProductRepository {
	return &OrderProductRepository{}
}

func (o *OrderProductRepository) Create(ctx context.Context, tx pgx.Tx, orderID, productID int) error {
	orderProductID, err := o.generateNextOrderProductID(ctx, tx)
	if err != nil {
		return err
	}
	orderProduct := &OrderProduct{ID: orderProductID, OrderID: orderID, ProductID: productID}

	sql, args := OrderProductStruct.InsertInto(OrderProductsTable, orderProduct).
		BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := tx.QueryRow(ctx, sql, args...)
	err = row.Scan()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	return nil
}

func (o *OrderProductRepository) generateNextOrderProductID(ctx context.Context, tx pgx.Tx) (int, error) {
	rows, err := tx.Query(ctx, fmt.Sprintf("SELECT nextval('%s')", "order_products_sequence"))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	return 0, fmt.Errorf("something was wrong. there is no next orderProduct id")
}

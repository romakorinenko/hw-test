package repository

import (
	"context"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
)

type OrderProduct struct {
	ID        int `db:"id" fieldtag:"pk" json:"id"`
	OrderId   int `db:"order_id" json:"orderId"`
	ProductId int `db:"product_id" json:"productId"`
}

type IOrderProductRepository interface {
	Create(ctx context.Context, tx pgx.Tx, orderId, productId int) error
}

type OrderProductRepository struct {
}

var OrderProductStruct = sqlbuilder.NewStruct(new(OrderProduct))

func NewOrderProductRepository() IOrderProductRepository {
	return &OrderProductRepository{}
}

func (o *OrderProductRepository) Create(ctx context.Context, tx pgx.Tx, orderId, productId int) error {
	orderProductId, err := o.generateNextOrderProductID(ctx, tx)
	if err != nil {
		return err
	}
	orderProduct := &OrderProduct{ID: orderProductId, OrderId: orderId, ProductId: productId}

	sql, args := OrderProductStruct.InsertInto("order_products", orderProduct).
		BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := tx.QueryRow(ctx, sql, args...)
	_ = row.Scan()

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

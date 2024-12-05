package repository

import (
	"context"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Order struct {
	ID          int       `db:"id" fieldtag:"pk" json:"id"`
	UserId      int       `db:"user_id" json:"userId"`
	OrderDate   time.Time `db:"order_date" json:"orderDate"`
	TotalAmount float32   `db:"total_amount" json:"totalAmount"`
	ProductIds  []int     `db:"-" json:"productIds"`
}

type IOrderRepository interface {
	Create(ctx context.Context, order *Order) (*Order, error)
	DeleteById(ctx context.Context, orderId int) error
	GetByUserId(ctx context.Context, userId int) ([]Order, error)
	GetByUserEmail(ctx context.Context, userEmail string) ([]Order, error)
}

type OrderRepository struct {
	dbPool                 *pgxpool.Pool
	OrderProductRepository IOrderProductRepository
}

var OrderStruct = sqlbuilder.NewStruct(new(Order))

func NewOrderRepository(dbPool *pgxpool.Pool) IOrderRepository {
	return &OrderRepository{dbPool: dbPool, OrderProductRepository: NewOrderProductRepository()}
}

func (o *OrderRepository) Create(ctx context.Context, order *Order) (*Order, error) {
	conn, err := o.dbPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	orderId, err := o.generateNextOrderId(ctx, tx)
	if err != nil {
		return nil, err
	}
	order.ID = orderId
	order.OrderDate = time.Now()

	sql, args := OrderStruct.InsertInto("orders", order).
		BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := tx.QueryRow(ctx, sql, args...)
	_ = row.Scan()

	for _, productId := range order.ProductIds {
		err := o.OrderProductRepository.Create(ctx, tx, int(orderId), productId)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *OrderRepository) DeleteById(ctx context.Context, orderId int) error {
	tx, err := o.dbPool.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadWrite})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	deleteBuilder := OrderStruct.DeleteFrom("orders")
	sql, args := deleteBuilder.Where(deleteBuilder.Equal("id", orderId)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("cannot delete order from db: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderRepository) GetByUserId(ctx context.Context, userId int) ([]Order, error) {
	selectBuilder := OrderStruct.SelectFrom("orders")
	sql, args := selectBuilder.Where(selectBuilder.Equal("user_id", userId)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := o.dbPool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Order, 0)
	for rows.Next() {
		var order Order
		err := rows.Scan(OrderStruct.Addr(&order)...)
		if err != nil {
			return nil, err
		}
		res = append(res, order)
	}

	return res, nil
}

func (o *OrderRepository) GetByUserEmail(ctx context.Context, userEmail string) ([]Order, error) {
	selectBuilder := OrderStruct.SelectFrom("orders")
	sql, args := selectBuilder.JoinWithOption(sqlbuilder.LeftJoin, "users", "orders.user_id = users.id").
		Where(selectBuilder.Equal("users.email", userEmail)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := o.dbPool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Order, 0)
	for rows.Next() {
		var order Order
		err := rows.Scan(OrderStruct.Addr(&order)...)
		if err != nil {
			return nil, err
		}
		res = append(res, order)
	}

	return res, nil
}

func (o *OrderRepository) generateNextOrderId(ctx context.Context, tx pgx.Tx) (int, error) {
	rows, err := tx.Query(ctx, fmt.Sprintf("SELECT nextval('%s')", "orders_sequence"))

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
	return 0, fmt.Errorf("something was wrong. there is no next order id")
}

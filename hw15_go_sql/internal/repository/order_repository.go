package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Order struct {
	ID          int       `db:"id" fieldtag:"pk" json:"id"`
	UserId      int       `db:"user_id" json:"userId"`
	OrderDate   time.Time `db:"order_date" json:"orderDate"`
	TotalAmount float32   `db:"total_amount" json:"totalAmount"`
	ProductIds  []int     `db:"-" json:"productIds"`
}

type UserStatistics struct {
	Name         string  `db:"name" json:"userName"`
	TotalOrders  int     `db:"total_orders" json:"totalOrders"`
	TotalAmount  float32 `db:"total_amount" json:"totalAmount"`
	AveragePrice float32 `db:"avg_price" json:"averagePrice"`
}

type IOrderRepository interface {
	Create(ctx context.Context, order *Order) (*Order, error)
	DeleteById(ctx context.Context, orderId int) error
	GetByUserID(ctx context.Context, userId int) ([]Order, error)
	GetByUserEmail(ctx context.Context, userEmail string) ([]Order, error)
	GetStatisticsByID(ctx context.Context, userId int) (*UserStatistics, error)
}

type OrderRepository struct {
	dbPool                 *pgxpool.Pool
	OrderProductRepository IOrderProductRepository
}

const ordersTable = "orders"

var OrderStruct = sqlbuilder.NewStruct(new(Order))
var UserStatisticsStruct = sqlbuilder.NewStruct(new(UserStatistics))

func NewOrderRepository(dbPool *pgxpool.Pool) *OrderRepository {
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

	orderId, err := o.generateNextOrderID(ctx, tx)
	if err != nil {
		return nil, err
	}
	order.ID = orderId
	order.OrderDate = time.Now()

	sql, args := OrderStruct.InsertInto(ordersTable, order).
		BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := tx.QueryRow(ctx, sql, args...)
	_ = row.Scan()

	for _, productId := range order.ProductIds {
		err := o.OrderProductRepository.Create(ctx, tx, orderId, productId)
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

	deleteBuilder := OrderStruct.DeleteFrom(ordersTable)
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

func (o *OrderRepository) GetByUserID(ctx context.Context, userId int) ([]Order, error) {
	selectBuilder := OrderStruct.SelectFrom(ordersTable)
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
	selectBuilder := OrderStruct.SelectFrom(ordersTable)
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

func (o *OrderRepository) GetStatisticsByID(ctx context.Context, userId int) (*UserStatistics, error) {
	selectBuilder := sqlbuilder.NewSelectBuilder()
	sql, args := selectBuilder.Select("users.name", "COUNT(DISTINCT orders.id) AS total_orders", "SUM(products.price) AS total_amount", "AVG(products.price) AS avg_price").
		From(ordersTable).
		Join("order_products", "orders.id = order_products.order_id").
		Join("products", "order_products.product_id = products.id").
		Join("users", "orders.user_id = users.id").
		Where(selectBuilder.Equal("users.id", userId)).
		GroupBy("users.id", "users.name").
		BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := o.dbPool.QueryRow(ctx, sql, args...)

	var userStatistics UserStatistics
	err := row.Scan(UserStatisticsStruct.Addr(&userStatistics)...)
	if err != nil {
		return nil, err
	}
	return &userStatistics, nil
}

func (o *OrderRepository) generateNextOrderID(ctx context.Context, tx pgx.Tx) (int, error) {
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
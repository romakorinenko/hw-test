package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	ID    int     `db:"id" fieldtag:"pk" json:"id"`
	Name  string  `db:"name" json:"name"`
	Price float32 `db:"price" json:"price"`
}

type IProductRepository interface {
	Create(ctx context.Context, product *Product) (*Product, error)
	GetByID(ctx context.Context, productID int) (*Product, error)
	GetAll(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, product *Product) (*Product, error)
	DeleteByID(ctx context.Context, productID int) error
}

type ProductRepository struct {
	dbPool *pgxpool.Pool
}

func NewProductRepository(dbPool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{dbPool: dbPool}
}

const productTable = "products"

var ProductStruct = sqlbuilder.NewStruct(new(Product))

func (p *ProductRepository) Create(ctx context.Context, product *Product) (*Product, error) {
	productID, err := p.generateNextProductID(ctx)
	if err != nil {
		return nil, err
	}
	product.ID = productID

	sql, args := ProductStruct.InsertInto(productTable, product).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	row := p.dbPool.QueryRow(ctx, sql, args...)
	rowScanErr := row.Scan()
	if rowScanErr != nil && !errors.Is(rowScanErr, pgx.ErrNoRows) {
		return nil, rowScanErr
	}

	return product, nil
}

func (p *ProductRepository) GetByID(ctx context.Context, productID int) (*Product, error) {
	selectBuilder := ProductStruct.SelectFrom(productTable)
	sql, args := selectBuilder.Where(selectBuilder.Equal("id", productID)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	row := p.dbPool.QueryRow(ctx, sql, args...)

	var product Product
	rowScanErr := row.Scan(ProductStruct.Addr(&product)...)
	if rowScanErr != nil {
		return nil, rowScanErr
	}

	return &product, nil
}

func (p *ProductRepository) GetAll(ctx context.Context) ([]Product, error) {
	sql, _ := ProductStruct.SelectFrom(productTable).
		OrderBy("id").
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := p.dbPool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Product, 0)
	for rows.Next() {
		var product Product
		rowScanErr := rows.Scan(ProductStruct.Addr(&product)...)
		if rowScanErr != nil {
			return nil, rowScanErr
		}
		res = append(res, product)
	}

	return res, nil
}

func (p *ProductRepository) Update(ctx context.Context, product *Product) (*Product, error) {
	updateBuilder := sqlbuilder.NewUpdateBuilder()
	sql, args := updateBuilder.Update(productTable).
		Set(
			updateBuilder.Assign("name", product.Name),
			updateBuilder.Assign("price", product.Price),
		).
		Where(updateBuilder.Equal("id", product.ID)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := p.dbPool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot update product: %w", err)
	}

	return product, nil
}

func (p *ProductRepository) DeleteByID(ctx context.Context, productID int) error {
	deleteBuilder := ProductStruct.DeleteFrom(productTable)
	sql, args := deleteBuilder.Where(deleteBuilder.Equal("id", productID)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := p.dbPool.Exec(ctx, sql, args...)
	return err
}

func (p *ProductRepository) generateNextProductID(ctx context.Context) (int, error) {
	rows, err := p.dbPool.Query(ctx, fmt.Sprintf("SELECT nextval('%s')", "products_sequence"))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		rowScanErr := rows.Scan(&id)
		if rowScanErr != nil {
			return 0, rowScanErr
		}
		return id, nil
	}
	return 0, fmt.Errorf("something was wrong. there is no next product id")
}

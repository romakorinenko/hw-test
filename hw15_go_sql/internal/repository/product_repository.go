package repository

import (
	"context"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	ID    int     `db:"id" fieldtag:"pk" json:"id"`
	Name  string  `db:"name" json:"name"`
	Price float32 `db:"price" json:"price"`
}

type IProductRepository interface {
	Create(ctx context.Context, product *Product) (*Product, error)
	GetByID(ctx context.Context, productId int) (*Product, error)
	GetAll(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, product *Product) (*Product, error)
	DeleteByID(ctx context.Context, productId int) error
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
	productId, err := p.generateNextProductID(ctx)
	if err != nil {
		return nil, err
	}
	product.ID = productId

	sql, args := ProductStruct.InsertInto(productTable, product).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_ = p.dbPool.QueryRow(ctx, sql, args...)

	return product, nil
}

func (p *ProductRepository) GetByID(ctx context.Context, productId int) (*Product, error) {
	selectBuilder := ProductStruct.SelectFrom(productTable)
	sql, args := selectBuilder.Where(selectBuilder.Equal("id", productId)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	row := p.dbPool.QueryRow(ctx, sql, args...)

	var product Product
	err := row.Scan(ProductStruct.Addr(&product)...)
	if err != nil {
		return nil, err
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
		err := rows.Scan(ProductStruct.Addr(&product)...)
		if err != nil {
			return nil, err
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

func (p *ProductRepository) DeleteByID(ctx context.Context, productId int) error {
	deleteBuilder := ProductStruct.DeleteFrom(productTable)
	sql, args := deleteBuilder.Where(deleteBuilder.Equal("id", productId)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := p.dbPool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("cannot delete product from db: %w", err)
	}

	return nil
}

func (p *ProductRepository) generateNextProductID(ctx context.Context) (int, error) {
	rows, err := p.dbPool.Query(ctx, fmt.Sprintf("SELECT nextval('%s')", "products_sequence"))

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
	return 0, fmt.Errorf("something was wrong. there is no next product id")
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"product_service/product"
)

type ProductRepository interface {
	Create(ctx context.Context, tx *sql.Tx, payload *product.ProductPayload) error
	GetProductByID(ctx context.Context, tx *sql.Tx, ID int) (product.Product, error)
	GetAllProduct(ctx context.Context, tx *sql.Tx) ([]*product.Product, error)
	UpdateProduct(ctx context.Context, tx *sql.Tx, payload *product.Product) error
	DeleteProduct(ctx context.Context, tx *sql.Tx, ID int) error
}

type ProductRepositoryImpl struct{}

func NewProductRepositoryImpl() *ProductRepositoryImpl {
	return &ProductRepositoryImpl{}
}

func (u *ProductRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, payload *product.ProductPayload) error {
	SQL := "insert into product(name, description, price, stock) values (?, ?, ?, ?)"
	if _, err := tx.ExecContext(ctx, SQL, payload.Name, payload.Description, payload.Price, payload.Stock); err != nil {
		return err
	}

	return nil
}

func (u *ProductRepositoryImpl) GetProductByID(ctx context.Context, tx *sql.Tx, ID int) (*product.Product, error) {
	SQL := "select id, name, description, price, stock, created_at, updated_at from product where id = ?"
	rows, err := tx.QueryContext(ctx, SQL, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productResponse := &product.Product{}
	if rows.Next() {
		if err := rows.Scan(&productResponse.Id, &productResponse.Name, &productResponse.Description, &productResponse.Price, &productResponse.Stock, &productResponse.CreatedAt, &productResponse.UpdatedAt); err != nil {
			return nil, err
		}
		return productResponse, nil
	} else {
		return productResponse, errors.New("product is not found")
	}
}

func (u *ProductRepositoryImpl) GetAllProduct(ctx context.Context, tx *sql.Tx) ([]*product.Product, error) {
	SQL := "SELECT id, name, description, price, stock, created_at, updated_at FROM product"
	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*product.Product
	for rows.Next() {
		productItem := &product.Product{}
		if err := rows.Scan(&productItem.Id, &productItem.Name, &productItem.Description, &productItem.Price, &productItem.Stock, &productItem.CreatedAt, &productItem.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, productItem)
	}

	return products, nil
}

func (u *ProductRepositoryImpl) UpdateProduct(ctx context.Context, tx *sql.Tx, payload *product.Product) error {
	SQL := "update product set name = ?, description = ?, price = ?, stock = ?, updated_at = ? where id = ?"
	if _, err := tx.ExecContext(ctx, SQL, payload.Name, payload.Description, payload.Price, payload.Stock, payload, payload.UpdatedAt, payload.Id); err != nil {
		return err
	}
	return nil
}

func (repository *ProductRepositoryImpl) DeleteProduct(ctx context.Context, tx *sql.Tx, ID int) error {
	SQL := "delete from product where id = ?"
	if _, err := tx.ExecContext(ctx, SQL, ID); err != nil {
		return err
	}
	return nil
}

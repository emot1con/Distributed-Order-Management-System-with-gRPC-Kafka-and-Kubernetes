package repository

import (
	"context"
	"database/sql"
	"errors"
	"product_service/proto"
)

type ProductRepository interface {
	Create(ctx context.Context, tx *sql.Tx, payload *proto.ProductPayload) error
	GetProductByID(ctx context.Context, tx *sql.Tx, ID int) (*proto.Product, error)
	GetAllProduct(ctx context.Context, db *sql.DB, offset int) ([]*proto.Product, int, int, error)
	UpdateProduct(ctx context.Context, tx *sql.Tx, payload *proto.Product) error
	DeleteProduct(ctx context.Context, tx *sql.Tx, ID int) error
}

type ProductRepositoryImpl struct{}

func NewProductRepositoryImpl() *ProductRepositoryImpl {
	return &ProductRepositoryImpl{}
}

func (u *ProductRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, payload *proto.ProductPayload) error {
	SQL := "INSERT INTO products(name, description, price, stock) VALUES ($1, $2, $3, $4)"
	if _, err := tx.ExecContext(ctx, SQL, payload.Name, payload.Description, payload.Price, payload.Stock); err != nil {
		return err
	}
	return nil
}

func (u *ProductRepositoryImpl) GetProductByID(ctx context.Context, tx *sql.Tx, ID int) (*proto.Product, error) {
	SQL := "SELECT id, name, description, price, stock, created_at, updated_at FROM products WHERE id = $1"
	rows := tx.QueryRowContext(ctx, SQL, ID)

	productResponse := &proto.Product{}
	if err := rows.Scan(
		&productResponse.Id,
		&productResponse.Name,
		&productResponse.Description,
		&productResponse.Price,
		&productResponse.Stock,
		&productResponse.CreatedAt,
		&productResponse.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return productResponse, nil
}

func (u *ProductRepositoryImpl) GetAllProduct(ctx context.Context, db *sql.DB, offset int) ([]*proto.Product, int, int, error) {
	SQL := "SELECT id, name, description, price, stock, created_at, updated_at FROM products ORDER BY id ASC LIMIT 15 OFFSET $1"
	rows, err := db.QueryContext(ctx, SQL, offset)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var products []*proto.Product
	for rows.Next() {
		productItem := &proto.Product{}
		if err := rows.Scan(&productItem.Id, &productItem.Name, &productItem.Description, &productItem.Price, &productItem.Stock, &productItem.CreatedAt, &productItem.UpdatedAt); err != nil {
			return nil, 0, 0, err
		}
		products = append(products, productItem)
	}

	var totalProduct int
	if err := db.QueryRow("SELECT COUNT(id) FROM products").Scan(&totalProduct); err != nil {
		return nil, 0, 0, err
	}
	totalPage := (totalProduct + 15 - 1) / 15

	return products, totalProduct, totalPage, nil
}

func (u *ProductRepositoryImpl) UpdateProduct(ctx context.Context, tx *sql.Tx, payload *proto.Product) error {
	SQL := `UPDATE products
        SET name = COALESCE(NULLIF($1, ''), name),
            description = COALESCE(NULLIF($2, ''), description),
            price = $3,
            stock = $4,
            updated_at = NOW()
        WHERE id = $5`
	if _, err := tx.ExecContext(ctx, SQL, payload.Name, payload.Description, payload.Price, payload.Stock, payload.Id); err != nil {
		return err
	}

	return nil
}

func (repository *ProductRepositoryImpl) DeleteProduct(ctx context.Context, tx *sql.Tx, ID int) error {
	SQL := "DELETE FROM products WHERE id = $1"
	result, err := tx.ExecContext(ctx, SQL, ID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("product not found")
	}

	return nil
}

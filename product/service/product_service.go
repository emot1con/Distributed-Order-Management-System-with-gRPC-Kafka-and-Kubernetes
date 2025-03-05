package service

import (
	"context"
	"database/sql"
	"product_service/helper"
	"product_service/product"
	"product_service/repository"
	"time"
)

type ProductService struct {
	repo repository.ProductRepository
	DB   *sql.DB
	ctx  context.Context
}

func NewProductService(ProductRepository repository.ProductRepository, DB *sql.DB, ctx context.Context) *ProductService {
	return &ProductService{
		repo: ProductRepository,
		DB:   DB,
		ctx:  ctx,
	}
}

func (u *ProductService) Create(payload *product.ProductRequest) (string, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return "", err
	}
	defer helper.CommitOrRollback(tx)

	if err := u.repo.Create(u.ctx, tx, payload.Payload); err != nil {
		return "", err
	}

	return "success add product", nil
}

func (u *ProductService) Update(payload *product.Product) (string, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return "", err
	}
	defer helper.CommitOrRollback(tx)

	productResult, err := u.repo.GetProductByID(u.ctx, tx, int(payload.Id))
	if err != nil {
		return "", err
	}

	productResult.Name = payload.Name
	productResult.Description = payload.Description
	productResult.Price = payload.Price
	productResult.Stock = payload.Stock
	productResult.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := u.repo.UpdateProduct(u.ctx, tx, payload); err != nil {
		return "", err
	}

	return "success update product", nil
}

func (u *ProductService) GetAll(page int) ([]*product.Product, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	offset := (page - 1) * 10
	products, err := u.repo.GetAllProduct(u.ctx, tx, offset)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (u *ProductService) Delete(ID int) (string, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return "", err
	}
	defer helper.CommitOrRollback(tx)

	if err := u.repo.DeleteProduct(u.ctx, tx, ID); err != nil {
		return "", err
	}

	return "success delete product", nil
}

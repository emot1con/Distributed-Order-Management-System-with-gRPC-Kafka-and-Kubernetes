package service

import (
	"context"
	"database/sql"
	"product_service/helper"
	"product_service/proto"
	"product_service/repository"

	"github.com/sirupsen/logrus"
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

func (u *ProductService) Create(payload *proto.ProductRequest) error {
	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	logrus.Info("calling create product service")
	if err := u.repo.Create(u.ctx, tx, payload.Payload); err != nil {
		return err
	}

	return nil
}

func (u *ProductService) GetUserByID(ID int) (*proto.Product, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	logrus.Info("getting product")
	return u.repo.GetProductByID(u.ctx, tx, ID)
}

func (u *ProductService) Update(payload *proto.Product) (*proto.Product, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	logrus.Info("getting product")
	productResult, err := u.repo.GetProductByID(u.ctx, tx, int(payload.Id))
	if err != nil {
		return nil, err
	}

	productResult.Name = payload.Name
	productResult.Description = payload.Description
	productResult.Price = payload.Price
	productResult.Stock = payload.Stock

	if err := u.repo.UpdateProduct(u.ctx, tx, payload); err != nil {
		return nil, err
	}
	logrus.Info("product updated")

	return productResult, nil
}

func (u *ProductService) GetAll(page int) ([]*proto.Product, int, int, int, error) {
	offset := (page - 1) * 10

	logrus.Info("getting all products")

	products, totalProducts, totalPage, err := u.repo.GetAllProduct(u.ctx, u.DB, offset)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	return products, totalProducts, totalPage, page, nil
}

func (u *ProductService) Delete(ID int) error {
	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	if err := u.repo.DeleteProduct(u.ctx, tx, ID); err != nil {
		return err
	}

	return nil
}

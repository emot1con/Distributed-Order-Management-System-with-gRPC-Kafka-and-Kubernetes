package repository

import (
	"broker/proto"
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductRepository interface {
	Create(payload *proto.ProductRequest) (*proto.Empty, error)
	GetProduct(payload *proto.GetProductRequest) (*proto.Product, error)
	ListProducts(offset *proto.Offset) (*proto.ProductList, error)
	UpdateProduct(payload *proto.Product) (*proto.Product, error)
	DeleteProduct(ID *proto.GetProductRequest) (*proto.Empty, error)
}

type ProductRepositoryImpl struct {
	client proto.ProductServiceClient
}

func NewProductRepository() *ProductRepositoryImpl {
	conn, err := grpc.NewClient("product_service:40001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Failed to connect: %v", err)
	}
	logrus.Info("Connected to product service")

	client := proto.NewProductServiceClient(conn)
	return &ProductRepositoryImpl{client: client}
}

func (u *ProductRepositoryImpl) Create(payload *proto.ProductRequest) (*proto.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.CreateProduct(ctx, payload)
}

func (u *ProductRepositoryImpl) GetProduct(payload *proto.GetProductRequest) (*proto.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.GetProduct(ctx, payload)
}

func (u *ProductRepositoryImpl) ListProducts(offset *proto.Offset) (*proto.ProductList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.ListProducts(ctx, offset)
}

func (u *ProductRepositoryImpl) UpdateProduct(payload *proto.Product) (*proto.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.UpdateProduct(ctx, payload)
}

func (u *ProductRepositoryImpl) DeleteProduct(ID *proto.GetProductRequest) (*proto.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.DeleteProduct(ctx, ID)
}

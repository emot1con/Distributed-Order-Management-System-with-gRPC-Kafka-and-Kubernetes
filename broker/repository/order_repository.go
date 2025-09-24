package repository

import (
	"broker/proto"
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderRepository interface {
	CreateOrder(*proto.CreateOrderRequest) (*proto.OrderResponse, error)
	GetOrder(*proto.GetOrderRequest) (*proto.OrderResponse, error)
}

type OrderRepositoryImpl struct {
	client proto.OrderServiceClient
}

func NewOrderRepository() *OrderRepositoryImpl {
	conn, err := grpc.NewClient("order-service:30001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Failed to connect: %v", err)
	}

	logrus.Info("Connected to product service")
	client := proto.NewOrderServiceClient(conn)

	return &OrderRepositoryImpl{client: client}
}

func (u *OrderRepositoryImpl) CreateOrder(payload *proto.CreateOrderRequest) (*proto.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logrus.Infof("Create the order and pass to order service with")
	return u.client.CreateOrder(ctx, payload)
}

func (u *OrderRepositoryImpl) GetOrder(payload *proto.GetOrderRequest) (*proto.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logrus.Infof("Get the order by ID and pass to order service with ID: %d", payload.OrderId)
	return u.client.GetOrder(ctx, payload)
}

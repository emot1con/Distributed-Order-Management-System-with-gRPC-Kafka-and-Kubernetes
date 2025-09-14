package repository

import (
	"context"
	"payment/proto"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderRepository interface {
	UpdateOrderStatus(ctx context.Context, req *proto.UpdateOrderStatusRequest) (*proto.EmptyOrder, error)
}

type OrderRepositoryImpl struct {
	client proto.OrderServiceClient
}

func NewOrderRepository() *OrderRepositoryImpl {
	conn, err := grpc.NewClient("order_service:30001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Failed to connect: %v", err)
	}

	logrus.Info("Connected to product service")
	client := proto.NewOrderServiceClient(conn)

	return &OrderRepositoryImpl{client: client}
}

func (u *OrderRepositoryImpl) UpdateOrderStatus(ctx context.Context, payload *proto.UpdateOrderStatusRequest) (*proto.EmptyOrder, error) {
	ctxRepo, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.UpdateOrderStatus(ctxRepo, payload)
}

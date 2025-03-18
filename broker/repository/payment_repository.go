package repository

import (
	"broker/proto"
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentRepository interface {
	PayOrder(*proto.CreatePaymentRequest) (*proto.OrderPayment, error)
	Transaction(*proto.PaymentTransaction) (*proto.EmptyPayment, error)
}

type PaymentRepositoryImpl struct {
	client proto.PaymentServiceClient
}

func NewPaymentRepository() *PaymentRepositoryImpl {
	conn, err := grpc.NewClient("payment_service:60001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Failed to connect: %v", err)
	}

	logrus.Info("Connected to payment service")
	client := proto.NewPaymentServiceClient(conn)

	return &PaymentRepositoryImpl{client: client}
}

func (u *PaymentRepositoryImpl) PayOrder(payload *proto.CreatePaymentRequest) (*proto.OrderPayment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.PayOrder(ctx, payload)
}

func (u *PaymentRepositoryImpl) Transaction(payload *proto.PaymentTransaction) (*proto.EmptyPayment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.Transaction(ctx, payload)
}

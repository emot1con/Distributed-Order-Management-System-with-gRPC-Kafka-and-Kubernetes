package service

import (
	"database/sql"
	"order/helper"
	"order/proto"
	"order/repository"
	"time"
)

type OrderService struct {
	DB            *sql.DB
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemsRepository
}

func NewOrderItemService(DB *sql.DB, orderRepo repository.OrderRepository, orderItemRepo repository.OrderItemsRepository) *OrderService {
	return &OrderService{
		DB:            DB,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
	}
}

func (u *OrderService) CreateOrder(payload *proto.CreateOrderRequest) (*proto.OrderResponse, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	orderID, err := u.orderRepo.CreateOrder(payload, payload.TotalPrice, tx)
	if err != nil {
		return nil, err
	}

	for _, v := range payload.Items {
		orderItemPayload := &proto.OrderItemRequest{
			OrderId:   int32(orderID),
			ProductId: v.ProductId,
			Quantity:  v.Quantity,
			Price:     v.Price,
		}
		err := u.orderItemRepo.CreateOrderItems(orderItemPayload, tx)
		if err != nil {
			return nil, err
		}
	}

	return &proto.OrderResponse{
		Order: &proto.Order{
			Id:         int32(orderID),
			TotalPrice: payload.TotalPrice,
			CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
			UpdatedAt:  time.Now().Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (u *OrderService) GetOrderByID(payload *proto.GetOrderRequest) (*proto.Order, error) {
	orderResponse, err := u.orderRepo.GetOrderByID(payload, u.DB)
	if err != nil {
		return nil, err
	}
	return orderResponse, nil
}

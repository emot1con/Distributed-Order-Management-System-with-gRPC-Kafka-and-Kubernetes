package repository

import (
	"database/sql"
	"errors"
	"order/proto"
)

type OrderRepository interface {
	CreateOrder(payload *proto.CreateOrderRequest, price float64, tx *sql.Tx) (int, error)
	GetOrderByID(payload *proto.GetOrderRequest, db *sql.DB) (*proto.Order, error)
}

type OrderRepositoryImpl struct{}

func NewOrderRepositoryImpl() *OrderRepositoryImpl {
	return &OrderRepositoryImpl{}
}

func (u *OrderRepositoryImpl) CreateOrder(payload *proto.CreateOrderRequest, price float64, tx *sql.Tx) (int, error) {
	SQL := "INSERT INTO orders(user_id, total_price) VALUES ($1, $2) RETURNING id"
	var orderID int
	err := tx.QueryRow(SQL, payload.UserId, price).Scan(&orderID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (u *OrderRepositoryImpl) GetOrderByID(payload *proto.GetOrderRequest, db *sql.DB) (*proto.Order, error) {
	SQL := "SELECT id, user_id, total_price, created_at, updated_at FROM orders WHERE id = $1"
	rows := db.QueryRow(SQL, payload.OrderId)

	orderResponse := &proto.Order{}
	if err := rows.Scan(
		&orderResponse.Id,
		&orderResponse.UserId,
		&orderResponse.TotalPrice,
		&orderResponse.CreatedAt,
		&orderResponse.UpdatedAt,
	); err != nil {
		return nil, errors.New("order not found")
	}
	return orderResponse, nil
}

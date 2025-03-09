package repository

import (
	"database/sql"
	"order/proto"
)

type OrderRepository interface {
	CreateOrder(payload *proto.CreateOrderRequest, tx *sql.Tx) error
	GetOrderByID(payload *proto.GetOrderRequest, db *sql.DB) (*proto.Order, error)
}

type OrderRepositoryImpl struct{}

func NewProductRepositoryImpl() *OrderRepositoryImpl {
	return &OrderRepositoryImpl{}
}

func (u *OrderRepositoryImpl) CreateOrder(payload *proto.CreateOrderRequest, tx *sql.Tx) error {
	SQL := "INSERT INTO orders(user_id, items) VALUES ($1, $2)"
	if _, err := tx.Exec(SQL, payload.UserId, payload.Items); err != nil {
		return err
	}
	return nil
}

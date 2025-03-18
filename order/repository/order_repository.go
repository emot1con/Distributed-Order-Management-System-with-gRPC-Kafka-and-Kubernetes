package repository

import (
	"database/sql"
	"errors"
	"order/proto"
	"time"
)

type OrderRepository interface {
	CreateOrder(payload *proto.CreateOrderRequest, price float64, tx *sql.Tx) (int, error)
	GetOrderByID(payload *proto.GetOrderRequest, db *sql.DB) (*proto.Order, error)
	UpdateOrderStatus(status string, orderID int, tx *sql.Tx) error
}

type OrderRepositoryImpl struct{}

func NewOrderRepositoryImpl() *OrderRepositoryImpl {
	return &OrderRepositoryImpl{}
}

func (u *OrderRepositoryImpl) CreateOrder(payload *proto.CreateOrderRequest, price float64, tx *sql.Tx) (int, error) {
	SQL := "INSERT INTO orders(user_id, total_price, status) VALUES ($1, $2, $3) RETURNING id"
	var orderID int
	err := tx.QueryRow(SQL, payload.UserId, price, "pending").Scan(&orderID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (u *OrderRepositoryImpl) GetOrderByID(payload *proto.GetOrderRequest, db *sql.DB) (*proto.Order, error) {
	SQL := "SELECT id, user_id, status, total_price, created_at, updated_at FROM orders WHERE id = $1"
	rows := db.QueryRow(SQL, payload.OrderId)

	orderResponse := &proto.Order{}
	if err := rows.Scan(
		&orderResponse.Id,
		&orderResponse.UserId,
		&orderResponse.Status,
		&orderResponse.TotalPrice,
		&orderResponse.CreatedAt,
		&orderResponse.UpdatedAt,
	); err != nil {
		return nil, errors.New("order not found")
	}
	return orderResponse, nil
}

func (u *OrderRepositoryImpl) UpdateOrderStatus(status string, orderID int, tx *sql.Tx) error {
	loc := time.FixedZone("WIB", 7*60*60)
	SQL := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`
	now := time.Now().In(loc)
	if _, err := tx.Exec(SQL, status, now, orderID); err != nil {
		return err
	}
	return nil
}

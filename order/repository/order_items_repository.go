package repository

import (
	"database/sql"
	"order/proto"
)

type OrderItemsRepository interface {
	CreateOrderItems(payload *proto.OrderItemRequest, tx *sql.Tx) error
	// GetOrderItemsByOrderID(payload *proto.GetOrderItemRequest, db *sql.DB) ([]*proto.OrderItem, error)
	// UpdateOrderItems(payload *proto.OrderItemRequest, tx *sql.Tx) error
	DeleteOrderItems(payload *proto.GetOrderItemRequest, tx *sql.Tx) error
}

type OrderItemsRepositoryImpl struct{}

func NewOrderItemsRepositoryImpl() *OrderItemsRepositoryImpl {
	return &OrderItemsRepositoryImpl{}
}

func (u *OrderItemsRepositoryImpl) CreateOrderItems(payload *proto.OrderItemRequest, tx *sql.Tx) error {
	SQL := "INSERT INTO order_items(order_id,product_id, quantity, price) VALUES ($1, $2, $3)"
	if _, err := tx.Exec(SQL, payload.OrderId, payload.ProductId, payload.Quantity, payload.Price); err != nil {
		return err
	}

	return nil
}

func (u *OrderItemsRepositoryImpl) DeleteOrderItems(payload *proto.GetOrderItemRequest, tx *sql.Tx) error {
	SQL := "DELETE FROM order_items WHERE id = $1"
	if _, err := tx.Exec(SQL, payload.OrderItemId); err != nil {
		return err
	}

	return nil
}

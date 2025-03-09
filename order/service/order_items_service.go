package service

// import (
// 	"database/sql"
// 	"order/helper"
// 	"order/proto"
// 	"order/repository"
// )

// type OrderItemService struct {
// 	DB            *sql.DB
// 	orderItemRepo repository.OrderItemsRepository
// }

// func NewOrderItemService(DB *sql.DB, orderItemRepo repository.OrderItemsRepository) *OrderItemService {
// 	return &OrderItemService{
// 		DB:            DB,
// 		orderItemRepo: orderItemRepo,
// 	}
// }

// func (u *OrderItemService) Create(payload *proto.OrderItemRequest) error {
// 	tx, err := u.DB.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer helper.CommitOrRollback(tx)

// 	if err := u.orderItemRepo.CreateOrderItems(payload, tx); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (u *OrderItemService) Delete(payload *proto.GetOrderItemRequest) error {
// 	tx, err := u.DB.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer helper.CommitOrRollback(tx)

// 	if err := u.orderItemRepo.DeleteOrderItems(payload, tx); err != nil {
// 		return err
// 	}
// 	return nil
// }

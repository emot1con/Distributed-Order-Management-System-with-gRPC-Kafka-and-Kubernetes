package service

import (
	"context"
	"database/sql"
	"payment/repository"
)

type PaymentService struct {
	repo repository.PaymentRepository
	DB   *sql.DB
	ctx  context.Context
}

func NewPaymentService(repo repository.PaymentRepository, DB *sql.DB, ctx context.Context) *PaymentService {
	return &PaymentService{
		repo: repo,
		DB:   DB,
		ctx:  ctx,
	}
}

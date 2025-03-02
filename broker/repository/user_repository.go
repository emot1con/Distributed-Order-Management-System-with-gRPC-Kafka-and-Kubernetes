package repository

import (
	"broker/usergrpc"
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserRepository interface {
	Register(*usergrpc.RegisterRequest) error
	Login(*usergrpc.LoginRequest) (*usergrpc.TokenResponse, error)
	RefreshToken(*usergrpc.RefreshTokenRequest) (*usergrpc.TokenResponse, error)
}

type UserRepositoryImpl struct {
	usergrpc.UnimplementedAuthServiceServer
}

func NewUserRepository() *UserRepositoryImpl {
	return &UserRepositoryImpl{
		UnimplementedAuthServiceServer: usergrpc.UnimplementedAuthServiceServer{},
	}
}

func (u *UserRepositoryImpl) Register(payload *usergrpc.RegisterRequest) error {
	conn, err := grpc.NewClient("user_service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logrus.Error("Failed to connect: ", err)
		return err
	}
	defer conn.Close()

	c := usergrpc.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logrus.Info("Connected to user service")

	_, err = c.Register(ctx, &usergrpc.RegisterRequest{
		Payload: &usergrpc.RegisterPayload{
			FullName: payload.Payload.FullName,
			Email:    payload.Payload.Email,
			Password: payload.Payload.Password,
		},
	})
	if err != nil {
		logrus.Error("Failed to register: ", err)
		return err
	}

	logrus.Info("Registered successfully")
	return nil
}

func (u *UserRepositoryImpl) Login(payload *usergrpc.LoginRequest) (*usergrpc.TokenResponse, error) {
	conn, err := grpc.NewClient("user_service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Error("Failed to connect: ", err)
		return nil, err
	}
	defer conn.Close()

	c := usergrpc.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logrus.Info("Connected to user service")

	response, err := c.Login(ctx, &usergrpc.LoginRequest{
		Payload: &usergrpc.LoginPayload{
			Email:    payload.Payload.Email,
			Password: payload.Payload.Password,
		},
	})
	if err != nil {
		logrus.Error("Failed to register: ", err)
		return nil, err
	}

	return response, nil
}

func (u *UserRepositoryImpl) RefreshToken(payload *usergrpc.RefreshTokenRequest) (*usergrpc.TokenResponse, error) {
	conn, err := grpc.NewClient("user_service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Error("Failed to connect: ", err)
		return nil, err
	}
	defer conn.Close()

	c := usergrpc.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logrus.Info("Connected to user service")

	response, err := c.RefreshToken(ctx, &usergrpc.RefreshTokenRequest{
		Payload: &usergrpc.RefreshTokenPayload{
			RefreshToken: payload.Payload.RefreshToken,
		},
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

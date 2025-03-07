package grpc

import (
	"context"
	"net"
	"time"

	"user_service/cmd/db"
	"user_service/proto"
	"user_service/repository"
	"user_service/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCServer struct {
	service *service.UserService
	proto.UnimplementedAuthServiceServer
}

func NewUserGRPCServer(service *service.UserService) *UserGRPCServer {
	return &UserGRPCServer{
		service: service,
	}
}

func (u *UserGRPCServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.EmptyResponse, error) {
	if err := u.service.Register(req.Payload); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed registering user: %v", err)
	}
	logrus.Info("user regisered successfully")

	return nil, nil
}

func (u *UserGRPCServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.TokenResponse, error) {
	token, refreshtoken, err := u.service.Login(req.Payload)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.TokenResponse{
		Message:               "success",
		Role:                  "user",
		Token:                 token,
		ExpiredAt:             time.Now().AddDate(0, 0, 1).String(),
		RefreshToken:          refreshtoken,
		RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0).String(),
	}, nil
}

func (u *UserGRPCServer) RefreshToken(ctx context.Context, req *proto.RefreshTokenRequest) (*proto.TokenResponse, error) {
	token, refreshtoken, err := u.service.RefreshToken(req.Payload.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.TokenResponse{
		Message:               "success",
		Role:                  "user",
		Token:                 token,
		ExpiredAt:             time.Now().AddDate(0, 0, 1).String(),
		RefreshToken:          refreshtoken,
		RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0).String(),
	}, nil
}

func GRPCListen() {
	DB, err := db.Connect()
	repo := repository.NewUserRepository(DB)
	service := service.NewUserService(repo)
	connection := NewUserGRPCServer(service)

	if err != nil {
		logrus.Fatalf("failed to connect to database: %v", err)
	}

	lis, err := net.Listen("tcp", ":50001")
	if err != nil {
		logrus.Fatalf("Failed to listen for gRPC: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterAuthServiceServer(srv, connection)
	logrus.Infof("gRPC Server started on port 50001")

	if err := srv.Serve(lis); err != nil {
		logrus.Fatalf("error when connect to gRPC Server  %v", err)
	}
}

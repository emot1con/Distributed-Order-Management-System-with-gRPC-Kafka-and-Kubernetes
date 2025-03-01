package grpcuser

import (
	"context"
	"net"
	"strings"
	"time"
	"user_service/service"
	"user_service/types"
	"user_service/user/usergrpc"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserGRPCServer struct {
	service service.UserService
	usergrpc.UnimplementedAuthServiceServer
}

func (u *UserGRPCServer) Register(ctx context.Context, req *usergrpc.RegisterRequest) (*usergrpc.EmptyResponse, error) {
	payload := req.GetPayload()

	dataLog := types.RegisterPayload{
		FullName: payload.FullName,
		Email:    payload.Email,
		Password: payload.Password,
	}

	if err := u.service.Register(dataLog); err != nil {
		return nil, status.Errorf(codes.Internal, "Gagal menyimpan user: %v", err)
	}

	return &usergrpc.EmptyResponse{}, nil
}

func (u *UserGRPCServer) Login(ctx context.Context, req *usergrpc.LoginRequest) (*usergrpc.TokenResponse, error) {
	payload := req.GetPayload()

	dataLog := types.LoginPayload{
		Email:    payload.Email,
		Password: payload.Password,
	}

	token, refreshtoken, err := u.service.Login(dataLog)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	return &usergrpc.TokenResponse{
		Message:               "success",
		Role:                  "user",
		Token:                 token,
		ExpiredAt:             time.Now().AddDate(0, 0, 1).String(),
		RefreshToken:          refreshtoken,
		RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0).String(),
	}, nil
}

func (u *UserGRPCServer) RefreshToken(ctx context.Context, req *usergrpc.RefreshTokenRequest) (*usergrpc.TokenResponse, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	authHeader, exists := md["authorization"]
	if !exists || len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
	}

	authToken := strings.TrimPrefix(authHeader[0], "Bearer ")
	if authToken == authHeader[0] {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token format")
	}

	token, refreshtoken, err := u.service.RefreshToken(authToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	return &usergrpc.TokenResponse{
		Message:               "success",
		Role:                  "user",
		Token:                 token,
		ExpiredAt:             time.Now().AddDate(0, 0, 1).String(),
		RefreshToken:          refreshtoken,
		RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0).String(),
	}, nil
}

func GRPCListen() {
	lis, err := net.Listen("tcp", ":50001")
	if err != nil {
		logrus.Fatalf("Failed to listen for gRPC: %v", err)
	}

	srv := grpc.NewServer()
	usergrpc.RegisterAuthServiceServer(srv, &UserGRPCServer{})
	logrus.Infof("gRPC Server started on port 50001")

	if err := srv.Serve(lis); err != nil {
		logrus.Fatalf("error when connect to gRPC Server  %v", err)
	}
}

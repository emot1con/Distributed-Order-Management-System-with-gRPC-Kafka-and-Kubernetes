package grpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"user_service/auth"
	"user_service/cmd/db"
	"user_service/proto"
	"user_service/repository"
	"user_service/service"
	"user_service/types"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
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

func (u *UserGRPCServer) GetUserByID(ctx context.Context, req *proto.GetUserRequest) (*proto.User, error) {
	user, err := u.service.GetUserByID(int(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.User{
		ID:        int32(user.ID),
		FullName:  user.FullName,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (u *UserGRPCServer) GoogleOauth(ctx context.Context, req *proto.EmptyRequest) (*proto.URLResponse, error) {
	url := auth.OauthGoogleConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return &proto.URLResponse{
		Url: url,
	}, nil
}

func (u *UserGRPCServer) FacebookOauth(ctx context.Context, req *proto.EmptyRequest) (*proto.URLResponse, error) {
	url := auth.OauthFacebookConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return &proto.URLResponse{
		Url: url}, nil
}

func (u *UserGRPCServer) GithubOauth(ctx context.Context, req *proto.EmptyRequest) (*proto.URLResponse, error) {
	url := auth.OauthGithubConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return &proto.URLResponse{
		Url: url,
	}, nil
}

func (u *UserGRPCServer) GoogleOauthCallback(ctx context.Context, req *proto.CodeRequest) (*proto.TokenResponse, error) {
	code := req.Code
	if code == "" {
		return nil, fmt.Errorf("missing code")
	}

	logrus.Info("handling Google OAuth authentication callback")

	token, err := auth.OauthGoogleConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := auth.OauthGoogleConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	oauthUser := &types.OAuthUserData{}
	if err := json.NewDecoder(resp.Body).Decode(oauthUser); err != nil {
		return nil, err
	}

	if oauthUser.Email == "" {
		return nil, fmt.Errorf("email not found or private")
	}

	oauthUser.Provider = "Google"

	logrus.Info("getting user by email from database")
	userResp, err := u.service.GetUserByEmail(oauthUser.Email)
	if err != nil && (err != sql.ErrNoRows && !errors.Is(err, gorm.ErrRecordNotFound)) {
		return nil, err
	} else if userResp != nil {
		userResp.Provider = "Google"
		userResp.ProviderId = oauthUser.ProviderID

		if err := u.service.UpdateUser(userResp); err != nil {
			return nil, err
		}

		jwtToken, err := auth.CreateFullJWTToken(int(userResp.ID))
		if err != nil {
			return nil, err
		}

		return jwtToken, nil
	}
	if err := u.service.Register(&proto.RegisterPayload{
		Email:      oauthUser.Email,
		Provider:   "Google",
		ProviderId: oauthUser.ProviderID,
		FullName:   oauthUser.Name,
	}); err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	user, err := u.service.GetUserByEmail(oauthUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}

	jwtToken, err := auth.CreateFullJWTToken(int(user.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT token: %v", err)
	}

	return jwtToken, nil
}

func (u *UserGRPCServer) FacebookOauthCallback(ctx context.Context, req *proto.CodeRequest) (*proto.TokenResponse, error) {
	code := req.Code
	if code == "" {
		return nil, fmt.Errorf("missing code")
	}

	logrus.Info("handling Facebook OAuth authentication callback")

	token, err := auth.OauthFacebookConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := auth.OauthFacebookConfig.Client(ctx, token)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	oauthUser := &types.OAuthUserData{}
	if err := json.NewDecoder(resp.Body).Decode(oauthUser); err != nil {
		return nil, err
	}

	if oauthUser.Email == "" {
		return nil, fmt.Errorf("email not found or private")
	}

	oauthUser.Provider = "Facebook"

	userResp, err := u.service.GetUserByEmail(oauthUser.Email)
	if err != nil && (err != sql.ErrNoRows && !errors.Is(err, gorm.ErrRecordNotFound)) {
		return nil, err
	} else if userResp != nil {
		userResp.Provider = "Facebook"
		userResp.ProviderId = oauthUser.ProviderID

		if err := u.service.UpdateUser(userResp); err != nil {
			return nil, err
		}

		jwtToken, err := auth.CreateFullJWTToken(int(userResp.ID))
		if err != nil {
			return nil, err
		}

		return jwtToken, nil
	}
	if err := u.service.Register(&proto.RegisterPayload{
		Email:      oauthUser.Email,
		Provider:   "Facebook",
		ProviderId: oauthUser.ProviderID,
		FullName:   oauthUser.Name,
	}); err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	user, err := u.service.GetUserByEmail(oauthUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}

	jwtToken, err := auth.CreateFullJWTToken(int(user.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT token: %v", err)
	}

	return jwtToken, nil
}

func (u *UserGRPCServer) GithubOauthCallback(ctx context.Context, req *proto.CodeRequest) (*proto.TokenResponse, error) {
	code := req.Code
	if code == "" {
		return nil, fmt.Errorf("missing code")
	}

	logrus.Info("handling Github OAuth authentication callback")

	token, err := auth.OauthGithubConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := auth.OauthGithubConfig.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	oauthUser := &types.OAuthUserData{}
	if err := json.NewDecoder(resp.Body).Decode(oauthUser); err != nil {
		return nil, err
	}

	if oauthUser.Email == "" {
		return nil, fmt.Errorf("email not found or private")
	}

	oauthUser.Provider = "Github"

	userResp, err := u.service.GetUserByEmail(oauthUser.Email)
	if err != nil && (err != sql.ErrNoRows && !errors.Is(err, gorm.ErrRecordNotFound)) {
		return nil, err
	} else if userResp != nil {
		userResp.Provider = "Github"
		userResp.ProviderId = oauthUser.ProviderID

		if err := u.service.UpdateUser(userResp); err != nil {
			return nil, err
		}

		jwtToken, err := auth.CreateFullJWTToken(int(userResp.ID))
		if err != nil {
			return nil, err
		}

		return jwtToken, nil
	}
	if err := u.service.Register(&proto.RegisterPayload{
		Email:      oauthUser.Email,
		Provider:   "Github",
		ProviderId: oauthUser.ProviderID,
		FullName:   oauthUser.Name,
	}); err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	user, err := u.service.GetUserByEmail(oauthUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}

	jwtToken, err := auth.CreateFullJWTToken(int(user.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT token: %v", err)
	}

	return jwtToken, nil
}

func GRPCListen() {
	auth.InitOauth()

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

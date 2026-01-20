package auth

import (
	"context"

	ssov1 "github.com/Adopten123/go-protobufcontract-sso/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyAppID  = 0
	emptyUserID = 0
	emptyString = ""
)

type AuthServer interface {
	Login(context context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(context context.Context, email string, password string) (userID int64, err error)
	IsAdmin(context context.Context, userID int64) (isAdmin bool, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	authServer AuthServer
}

func Register(gRPC *grpc.Server, authServer AuthServer) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{
		authServer: authServer,
	})
}

func (s *serverAPI) Login(context context.Context, request *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if err := validateLogin(request); err != nil {
		return nil, err
	}

	token, err := s.authServer.Login(context, request.GetEmail(), request.GetPassword(), int(request.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(context context.Context, request *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(request); err != nil {
		return nil, err
	}

	userID, err := s.authServer.RegisterNewUser(context, request.GetEmail(), request.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(context context.Context, request *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(request); err != nil {
		return nil, err
	}

	isAdmin, err := s.authServer.IsAdmin(context, request.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateLogin(request *ssov1.LoginRequest) error {
	if request.GetEmail() == emptyString {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if request.GetPassword() == emptyString {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if request.GetAppId() == emptyAppID {
		return status.Error(codes.InvalidArgument, "appId is required")
	}
	return nil
}

func validateRegister(request *ssov1.RegisterRequest) error {
	if request.GetEmail() == emptyString {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if request.GetPassword() == emptyString {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}

func validateIsAdmin(request *ssov1.IsAdminRequest) error {
	if request.GetUserId() == emptyUserID {
		return status.Error(codes.InvalidArgument, "userId is required")
	}
	return nil
}

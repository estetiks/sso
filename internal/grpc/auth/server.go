package authgrpc

import (
	"context"

	ssov1 "github.com/estetiks/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue = 0
	admin      = true
	notAdmin   = false
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)

	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponce, error) {

	if err := registerValidate(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())

	if err != nil {
		//TODO: ...create correct error
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponce{
		UserId: userID,
	}, nil

}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponce, error) {

	if err := loginValidate(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))

	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	role := isAdmin(token)

	return &ssov1.LoginResponce{
		Token:   token,
		IsAdmin: role,
	}, nil
}

func loginValidate(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "empty email")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "empty password")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "appID is required")
	}

	return nil
}

func registerValidate(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "empty email")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "empty password")
	}

	return nil
}

func isAdmin(token string) bool {
	//TODO: write code for indicate admin
	return notAdmin
}

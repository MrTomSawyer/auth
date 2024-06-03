package auth

import (
	"context"
	"errors"
	"fmt"
	ssov1 "github.com/MrTomSawyer/protos/gen/go/sso"
	"github.com/MrTomSawyer/sso/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "something went wrong")
	}
	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func (s *serverAPI) User(ctx context.Context, none *emptypb.Empty) (*ssov1.UserResponse, error) {
	user := ctx.Value("user").(jwt.MapClaims)
	fmt.Printf("TOKEN! : %s", user)
	return &ssov1.UserResponse{UserExists: true}, nil
}

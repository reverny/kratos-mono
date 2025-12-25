package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/reverny/kratos-mono/gen/go/api/user/v1"
	"github.com/reverny/kratos-mono/services/user/internal/biz"
)

type UserService struct {
	v1.UnimplementedUserServer

	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

func (s *UserService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.UserInfo, error) {
	return s.uc.CreateUser(ctx, req)
}

func (s *UserService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.UserInfo, error) {
	return s.uc.GetUser(ctx, req)
}

func (s *UserService) ListUsers(ctx context.Context, req *v1.ListUsersRequest) (*v1.ListUsersResponse, error) {
	return s.uc.ListUsers(ctx, req)
}

func (s *UserService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UserInfo, error) {
	return s.uc.UpdateUser(ctx, req)
}

func (s *UserService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	return s.uc.DeleteUser(ctx, req)
}

func (s *UserService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	return s.uc.Login(ctx, req)
}

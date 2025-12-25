package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/reverny/kratos-mono/gen/go/api/user/v1"
)

// UserRepo is a User repo.
type UserRepo interface {
	CreateUser(context.Context, *v1.CreateUserRequest) (*v1.UserInfo, error)
	GetUser(context.Context, string) (*v1.UserInfo, error)
	GetUserByUsername(context.Context, string) (*v1.UserInfo, string, error)
	ListUsers(context.Context, int32, int32, string, string) ([]*v1.UserInfo, int32, error)
	UpdateUser(context.Context, *v1.UpdateUserRequest) (*v1.UserInfo, error)
	DeleteUser(context.Context, string) error
}

// UserUsecase is a User usecase.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase new a User usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateUser creates a User.
func (uc *UserUsecase) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.UserInfo, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", req.Username)
	return uc.repo.CreateUser(ctx, req)
}

// GetUser gets a User by ID.
func (uc *UserUsecase) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.UserInfo, error) {
	uc.log.WithContext(ctx).Infof("GetUser: %v", req.Id)
	return uc.repo.GetUser(ctx, req.Id)
}

// ListUsers lists all Users.
func (uc *UserUsecase) ListUsers(ctx context.Context, req *v1.ListUsersRequest) (*v1.ListUsersResponse, error) {
	uc.log.WithContext(ctx).Infof("ListUsers: page=%d, page_size=%d", req.Page, req.PageSize)

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	users, total, err := uc.repo.ListUsers(ctx, req.Page, req.PageSize, req.Role, req.Status)
	if err != nil {
		return nil, err
	}

	return &v1.ListUsersResponse{
		Users: users,
		Total: total,
	}, nil
}

// UpdateUser updates a User.
func (uc *UserUsecase) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UserInfo, error) {
	uc.log.WithContext(ctx).Infof("UpdateUser: %v", req.Id)
	return uc.repo.UpdateUser(ctx, req)
}

// DeleteUser deletes a User.
func (uc *UserUsecase) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	uc.log.WithContext(ctx).Infof("DeleteUser: %v", req.Id)
	err := uc.repo.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Login authenticates a user.
func (uc *UserUsecase) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	uc.log.WithContext(ctx).Infof("Login: %v", req.Username)

	user, passwordHash, err := uc.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// TODO: Verify password with bcrypt
	_ = passwordHash

	// TODO: Generate JWT token
	token := "sample-jwt-token"

	return &v1.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

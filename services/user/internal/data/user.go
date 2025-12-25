package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"

	v1 "github.com/reverny/kratos-mono/gen/go/api/user/v1"
	"github.com/reverny/kratos-mono/services/user/internal/biz"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.UserInfo, error) {
	now := time.Now().Format(time.RFC3339)
	user := &v1.UserInfo{
		Id:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		FullName:  req.FullName,
		Role:      "user",
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// TODO: Hash password and save to database
	r.log.Infof("User created: %s", user.Id)

	return user, nil
}

func (r *userRepo) GetUser(ctx context.Context, id string) (*v1.UserInfo, error) {
	// TODO: Get from database
	user := &v1.UserInfo{
		Id:        id,
		Username:  "johndoe",
		Email:     "john@example.com",
		FullName:  "John Doe",
		AvatarUrl: "https://example.com/avatar.jpg",
		Role:      "user",
		Status:    "active",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	return user, nil
}

func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*v1.UserInfo, string, error) {
	// TODO: Get from database
	user := &v1.UserInfo{
		Id:        uuid.New().String(),
		Username:  username,
		Email:     "user@example.com",
		FullName:  "Sample User",
		Role:      "user",
		Status:    "active",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	passwordHash := "$2a$10$hashedpassword"

	return user, passwordHash, nil
}

func (r *userRepo) ListUsers(ctx context.Context, page, pageSize int32, role, status string) ([]*v1.UserInfo, int32, error) {
	// TODO: Get from database with pagination and filters
	users := []*v1.UserInfo{
		{
			Id:        "1",
			Username:  "admin",
			Email:     "admin@example.com",
			FullName:  "Admin User",
			Role:      "admin",
			Status:    "active",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
		{
			Id:        "2",
			Username:  "user1",
			Email:     "user1@example.com",
			FullName:  "User One",
			Role:      "user",
			Status:    "active",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
	}

	return users, int32(len(users)), nil
}

func (r *userRepo) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UserInfo, error) {
	// TODO: Update in database
	user := &v1.UserInfo{
		Id:        req.Id,
		Email:     req.Email,
		FullName:  req.FullName,
		AvatarUrl: req.AvatarUrl,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	r.log.Infof("User updated: %s", user.Id)

	return user, nil
}

func (r *userRepo) DeleteUser(ctx context.Context, id string) error {
	// TODO: Delete from database
	r.log.Infof("User deleted: %s", id)
	return nil
}

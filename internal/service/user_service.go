package service

import (
	"context"
	"errors"
	"fmt"

	"WalletApps/internal/model"      // 假设你有 model 层
	"WalletApps/internal/repository" // 假设你有 repository 层

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService 定义用户相关的服务接口
type UserService struct {
	userRepo *repository.UserRepository // 依赖注入用户数据访问层
}

// NewUserService 构造函数
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser 处理用户注册的核心业务逻辑
func (s *UserService) CreateUser(uid uuid.UUID, username, passwordHash, salt string) error {
	ctx := context.Background()

	existingUser, _ := s.userRepo.GetByUsername(ctx, username)
	if existingUser != nil {
		return errors.New("username already exists")
	}

	user := &model.User{
		ID:           uid,
		Name:         username,
		PasswordHash: passwordHash,
		Salt:         salt,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		fmt.Printf("DB Error: %v\n", err)
		return errors.New("failed to save user to database")
	}

	return nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

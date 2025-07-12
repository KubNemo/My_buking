package usecase

import (
	"context"
	"duking/internal/config"
	"duking/internal/models"
	"duking/internal/repository"
	"duking/internal/utils"
	"errors"
	"fmt"
	"strings"
	"time"
)

type UserService interface {
	Register(ctx context.Context, input models.RegisterInput) error
	Login(ctx context.Context, email, password string) (string, error)
	GetProfile(ctx context.Context, id uint) (*models.User, error)
	UpdateProfile(ctx context.Context, id uint, input models.UpdateInput) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

var cfg config.Config

func (s *userService) Register(ctx context.Context, input models.RegisterInput) error {
	if input.Email == "" {
		return errors.New("email not specified")
	}
	if !strings.Contains(input.Email, "@") {
		return errors.New("invalid email")
	}
	if input.Username == "" || input.Password == "" {
		return errors.New("not username and password")
	}

	exsistenUser, _ := s.repo.GetByEmail(ctx, input.Email)
	if exsistenUser != nil {
		return errors.New("a user with this email already exists")
	}
	hashPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		Username:   input.Username,
		Email:      input.Email,
		Password:   hashPassword,
		Role:       "user",
		CreatedAt:  time.Now(),
		IsVerified: false,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("email or password is empty")
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	if err := utils.ComparePassword(user.Password, password); err != nil {
		return "", errors.New("invalid password")
	}
	jwtManager := utils.NewJWTManager(cfg.SekretKey, time.Hour*2)
	token, err := jwtManager.Generate(user.UserID, user.Email, user.Role, user.Username)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return token, nil
}

func (s *userService) GetProfile(ctx context.Context, id uint) (*models.User, error) {
	// TODO: implement get profile logic
	return nil, nil
}

func (s *userService) UpdateProfile(ctx context.Context, id uint, input models.UpdateInput) error {
	// TODO: implement update profile logic
	return nil
}

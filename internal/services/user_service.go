package services

import (
	"fmt"
	"vasvault/internal/dto"
	"vasvault/internal/models"
	"vasvault/internal/repositories"
	"vasvault/pkg/utils"
	apperrors "vasvault/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	Register(request dto.RegisterRequest) (*dto.UserResponse, error)
	Login(request dto.LoginRequest) (*dto.UserResponse, error)
	GetUser(id uint) (*models.User, error)
	GetUserByID(id uint) (*dto.UserResponse, error)
	UpdateUser(id uint, request dto.UpdateProfileRequest) (*dto.UserResponse, error)
	Refresh(refreshToken string) (*dto.UserResponse, error)
}

type UserService struct {
	repository repositories.UserRepositoryInterface
}

func NewUserService(repo repositories.UserRepositoryInterface) UserServiceInterface {
	return &UserService{repository: repo}
}

func (s *UserService) Register(request dto.RegisterRequest) (*dto.UserResponse, error) {
	if _, err := s.repository.FindByEmail(request.Email); err == nil {
		return nil, apperrors.ErrEmailExists
	}
	if _, err := s.repository.FindByUsername(request.Username); err == nil {
		return nil, apperrors.ErrUsernameExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hash),
	}

	if err := s.repository.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	response := &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
	return response, nil
}

func (s *UserService) Login(request dto.LoginRequest) (*dto.UserResponse, error) {

	user, err := s.repository.FindByEmail(request.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	response := &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
	return response, nil
}

func (s *UserService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.repository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	response := &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
	return response, nil
}

func (s *UserService) UpdateUser(id uint, request dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.repository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if request.Email != "" {
		user.Email = request.Email
	}

	if request.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = string(hash)
	}

	if request.Username != "" {
		user.Username = request.Username
	}

	if err := s.repository.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
	return response, nil
}

func (s *UserService) GetUser(id uint) (*models.User, error) {
	user, err := s.repository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (s *UserService) Refresh(refreshToken string) (*dto.UserResponse, error) {
	token, err := utils.ValidateRefreshToken(refreshToken)

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	user, err := s.repository.FindByID(token.ID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	response := &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
	return response, nil
}

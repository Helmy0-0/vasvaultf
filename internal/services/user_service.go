package services

import (
	"fmt"
	"vasvault/internal/dto"
	"vasvault/internal/models"
	"vasvault/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.UserRepository
}

func NewAuthService(r repositories.UserRepository) UserService {
	return UserService{
		repository: r,
	}
}

func (s *UserService) Register(request dto.RegisterRequest) (*dto.UserResponse, error) {
	if _, err := s.repository.FindByEmail(request.Email); err == nil {
		return nil, fmt.Errorf("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{Email: request.Email, Password: string(hash)}
	if err := s.repository.Create(user); err != nil {
		return nil, err
	}

	response := &dto.UserResponse{ID: user.ID, Email: user.Email}
	return response, nil
}

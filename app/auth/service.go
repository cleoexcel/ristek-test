package auth

import (
	"errors"


	"github.com/cleoexcel/ristek-test/app/models"
)

type AuthService interface {
	GetAllUsers() ([]*models.User, error)
	Register(username, password string) (*models.User, error)
	Login(username, password string) (string, error)
	GetUser(username string) (*models.User, error)
}

type authService struct {
	repo Repository
}

func NewAuthService(repo Repository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) GetUser(username string) (*models.User, error) {
	return s.repo.GetUserByUsername(username)
}

func (s *authService) GetAllUsers() ([]*models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *authService) Register(username, password string) (*models.User, error) {
	user, err := s.repo.CreateUser(username, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(username, password string) (string, error) {
	err := s.repo.LoginUser(username, password)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}
	accessToken, err := GenerateAuthTokens(user) 
	if err != nil {
		return "", errors.New("failed to generate access token")
	}
	return accessToken, nil
}

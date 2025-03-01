package tryout

import (
	"github.com/cleoexcel/ristek-test/app/models"
)

type TryoutService interface {
	CreateTryout(title, description string, userId int, category string) (*models.Tryout, error)
	GetAllTryout(title string, category string, createdAt string) ([]*models.Tryout, error)
	GetDetailTryout(id int) (*models.Tryout, error)
	EditTryout(id int, title, description string, userId int) (*models.Tryout, error)
	DeleteTryoutById(id int) error
}

type tryoutService struct {
	repo Repository
}

func NewTryoutService(repo Repository) TryoutService {
	return &tryoutService{repo}
}

func (s *tryoutService) CreateTryout(title, description string, userId int, category string) (*models.Tryout, error) {
	tryout, err := s.repo.CreateTryout(title, description, userId, category)
	if err != nil {
		return nil, err
	}
	return tryout, nil
}

func (s *tryoutService) GetAllTryout(title string, category string, createdAt string) ([]*models.Tryout, error) {
	return s.repo.GetAllTryout(title, category, createdAt)
}

func (s *tryoutService) GetDetailTryout(id int) (*models.Tryout, error) {
	return s.repo.GetDetailTryout(id)
}

func (s *tryoutService) EditTryout(id int, title, description string, userId int) (*models.Tryout, error) {
	return s.repo.EditTryout(id, title, description, userId)
}

func (s *tryoutService) DeleteTryoutById(id int) error {
	return s.repo.DeleteTryoutById(id)
}

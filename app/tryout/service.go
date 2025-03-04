package tryout

import (
	"github.com/cleoexcel/ristek-test/app/models"
)

type TryoutService interface {
	CreateTryout(title, description string, userId int, category string) (*models.Tryout, error)
	GetAllTryout(title string, category string, createdAt string, userId int) ([]*models.Tryout, error)
	GetDetailTryoutByTryoutID(id int) (*models.Tryout, error)
	EditTryoutByTryoutID(id int, title, description string, userId int) (*models.Tryout, error)
	DeleteTryoutByTryoutID(id int) error
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

func (s *tryoutService) GetAllTryout(title string, category string, createdAt string, userId int) ([]*models.Tryout, error) {
	return s.repo.GetAllTryout(title, category, createdAt, userId)
}

func (s *tryoutService) GetDetailTryoutByTryoutID(id int) (*models.Tryout, error) {
	return s.repo.GetDetailTryoutByTryoutID(id)
}

func (s *tryoutService) EditTryoutByTryoutID(id int, title, description string, userId int) (*models.Tryout, error) {
	return s.repo.EditTryoutByTryoutID(id, title, description, userId)
}

func (s *tryoutService) DeleteTryoutByTryoutID(id int) error {
	return s.repo.DeleteTryoutByTryoutID(id)
}

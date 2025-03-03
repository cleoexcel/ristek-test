package answer

import (
	"github.com/cleoexcel/ristek-test/app/models"
)

type Service interface {
	CreateTruefalseAnswer(questionID int, expectAnswer bool) (*models.Truefalse, error)
	CreateShortanswerAnswer(questionID int, expectAnswer string) (*models.Shortanswer, error)
	GetAllAnswers() ([]*models.Truefalse, []*models.Shortanswer, error)
	EditTruefalseAnswer(id int, expectAnswer bool) (*models.Truefalse, error)
	EditShortanswerAnswer(id int, expectAnswer string) (*models.Shortanswer, error)
	DeleteAnswer(id int) error
	
}

type answerService struct {
	repo Repository
}

func NewAnswerService(repo Repository) Service {
	return &answerService{repo}
}

func (s *answerService) CreateTruefalseAnswer(questionID int, expectAnswer bool) (*models.Truefalse, error) {
	return s.repo.CreateTruefalseAnswer(questionID, expectAnswer)
}

func (s *answerService) CreateShortanswerAnswer(questionID int, expectAnswer string) (*models.Shortanswer, error) {
	return s.repo.CreateShortanswerAnswer(questionID, expectAnswer)
}

func (s *answerService) GetAllAnswers() ([]*models.Truefalse, []*models.Shortanswer, error) {
	return s.repo.GetAllAnswers()
}

func (s *answerService) EditTruefalseAnswer(id int, expectAnswer bool) (*models.Truefalse, error) {
	return s.repo.EditTruefalseAnswer(id, expectAnswer)
}

func (s *answerService) EditShortanswerAnswer(id int, expectAnswer string) (*models.Shortanswer, error) {
	return s.repo.EditShortanswerAnswer(id, expectAnswer)
}

func (s *answerService) DeleteAnswer(id int) error {
	return s.repo.DeleteAnswer(id)
}

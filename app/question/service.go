package question

import (
	"github.com/cleoexcel/ristek-test/app/models"
)

type QuestionService interface {
	CreateQuestion(content string, tryoutID int, questionType string, weight int) error
	GetAllQuestions(tryoutID int) ([]*models.Question, error)
	EditQuestion(id int, content string, questionType string, weight int) (*models.Question, error)
	DeleteQuestion(id int) error
}

type questionService struct {
	repo Repository
}

func NewQuestionService(repo Repository) QuestionService {
	return &questionService{repo}
}

func (s *questionService) CreateQuestion(content string, tryoutID int, questionType string, weight int) error {
	_, err := s.repo.CreateQuestion(content, tryoutID, questionType, weight)
	if err != nil {
		return err
	}
	return nil
}

func (s *questionService) GetAllQuestions(tryoutID int) ([]*models.Question, error) {
	return s.repo.GetAllQuestions(tryoutID)
}

func (s *questionService) EditQuestion(id int, content string, questionType string, weight int) (*models.Question, error) {
	return s.repo.EditQuestion(id, content, questionType, weight)
}

func (s *questionService) DeleteQuestion(id int) error {
	return s.repo.DeleteQuestion(id)
}

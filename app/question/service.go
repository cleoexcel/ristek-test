package question

import (
	"github.com/cleoexcel/ristek-test/app/answer"
	"github.com/cleoexcel/ristek-test/app/models"
)

type QuestionService struct {
	Repo          *QuestionRepository
	AnswerService *answer.AnswerService
}

func NewQuestionService(repo *QuestionRepository, answerService *answer.AnswerService) *QuestionService {
	return &QuestionService{Repo: repo, AnswerService: answerService}
}

func (s *QuestionService) GetAllQuestionsByTryoutID(tryoutID int) ([]*models.Question, error) {
	return s.Repo.GetAllQuestionsByTryoutID(tryoutID)
}

func (s *QuestionService) CreateQuestion(content string, tryoutID int, questionType string, weight int, expectAnswer interface{}) (*models.Question, error) {
	question, err := s.Repo.CreateQuestion(content, tryoutID, questionType, weight)
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (s *QuestionService) EditQuestionByQuestionID(id int, content string, weight int, expectAnswer interface{}) error {
	_, err := s.Repo.EditQuestionByQuestionID(id, content, weight)
	if err != nil {
		return err
	}
	_, err = s.AnswerService.UpdateAnswer(id, expectAnswer, []models.MultipleChoiceOption{})
	return err
}

func (s *QuestionService) DeleteQuestionByQuestionID(id int) error {
	return s.Repo.DeleteQuestionByQuestionID(id)
}

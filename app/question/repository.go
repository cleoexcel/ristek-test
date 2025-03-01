package question

import (
	"fmt"
	"github.com/cleoexcel/ristek-test/app/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllQuestions(tryoutID int) ([]*models.Question, error)
	CreateQuestion(number int, content string, tryoutID int, questionType string) (*models.Question, error)
	EditQuestion(id int, number int, content string, questionType string) (*models.Question, error)
	DeleteQuestion(id int) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{DB: db}
}

func (r *repository) GetAllQuestions(tryoutID int) ([]*models.Question, error) {
	var questions []*models.Question
	err := r.DB.Where("tryout_id = ?", tryoutID).Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *repository) CreateQuestion(number int, content string, tryoutID int, questionType string) (*models.Question, error) {
	var submission models.Submission
	err := r.DB.Where("tryout_id = ?", tryoutID).First(&submission).Error
	if err == nil {
		return nil, fmt.Errorf("tryout already has a submission, cannot add or edit questions")
	}

	question := &models.Question{
		Number:      number,
		Content:     content,
		TryoutID:    tryoutID,
		QuestionType: questionType,
	}
	if err := r.DB.Create(question).Error; err != nil {
		return nil, err
	}
	return question, nil
}

func (r *repository) EditQuestion(id int, number int, content string, questionType string) (*models.Question, error) {
	var question models.Question
	err := r.DB.First(&question, id).Error
	if err != nil {
		return nil, err
	}

	var submission models.Submission
	err = r.DB.Where("tryout_id = ?", question.TryoutID).First(&submission).Error
	if err == nil {
		return nil, fmt.Errorf("cannot edit question for tryout that already has a submission")
	}

	question.Number = number
	question.Content = content
	question.QuestionType = questionType
	if err := r.DB.Save(&question).Error; err != nil {
		return nil, err
	}

	return &question, nil
}

func (r *repository) DeleteQuestion(id int) error {
	var question models.Question
	err := r.DB.First(&question, id).Error
	if err != nil {
		return err
	}

	var submission models.Submission
	err = r.DB.Where("tryout_id = ?", question.TryoutID).First(&submission).Error
	if err == nil {
		return fmt.Errorf("cannot delete question for tryout that already has a submission")
	}

	err = r.DB.Delete(&question).Error
	if err != nil {
		return err
	}
	return nil
}

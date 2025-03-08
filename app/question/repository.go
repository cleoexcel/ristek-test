package question

import (
	"fmt"

	"github.com/cleoexcel/ristek-test/app/models"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	DB *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{DB: db}
}

func (r *QuestionRepository) GetAllQuestionsByTryoutID(tryoutID int) ([]*models.Question, error) {
	var questions []*models.Question
	err := r.DB.
	Preload("Tryout").
		Preload("ShortAnswer").
		Preload("TrueFalse").
		Preload("MultipleChoice").
		Preload("MultipleChoice.Options").
		Where("tryout_id = ?", tryoutID).
		Find(&questions).
		Error
	return questions, err
}

func (r *QuestionRepository) CreateQuestion(content string, tryoutID int, questionType string, weight int) (*models.Question, error) {
	var numberOfAttempt int64

 	r.DB.Model(&models.Submission{}).Where("tryout_id = ?", tryoutID).Count(&numberOfAttempt)

	if numberOfAttempt > 0 {
		return nil, fmt.Errorf("tryout already has a submission, cannot add or edit questions")
	}

	question := &models.Question{
		Number:       int(numberOfAttempt),
		Content:      content,
		TryoutID:     tryoutID,
		QuestionType: questionType,
		Weight:       weight,
	}
	if err := r.DB.Create(question).Error; err != nil {
		return nil, err
	}
	

	if err := r.DB.Preload("Tryout").First(&question, question.ID).Error; err != nil {
		return nil, err
	}

	return question, nil
}

func (r *QuestionRepository) EditQuestionByQuestionID(id int, content string, weight int) (*models.Question, error) {
	var question models.Question
	err := r.DB.
		Preload("Tryout").
		Preload("ShortAnswer").
		Preload("TrueFalse").
		Preload("MultipleChoice").
		Preload("MultipleChoice.Options").
		First(&question, id).Error
	if err != nil {
		return nil, err
	}

	var submission models.Submission
	err = r.DB.Where("tryout_id = ?", question.TryoutID).First(&submission).Error
	if err == nil {
		return nil, fmt.Errorf("cannot edit question for tryout that already has a submission")
	}

	question.Content = content
	question.Weight = weight
	if err := r.DB.Save(&question).Error; err != nil {
		return nil, err
	}

	return &question, nil
}

func (r *QuestionRepository) DeleteQuestionByQuestionID(id int) error {
	var question models.Question
	err := r.DB.Preload("Tryout").First(&question, id).Error 
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

func (r *QuestionRepository) GetQuestionByID(id int) (*models.Question, error) {
	var question models.Question
	err := r.DB.
		Preload("Tryout"). 
		Preload("ShortAnswer").
		Preload("TrueFalse").
		Preload("MultipleChoice").
		First(&question, id).
		Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}
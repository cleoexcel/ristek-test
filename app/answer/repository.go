package answer

import (
	"github.com/cleoexcel/ristek-test/app/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateTruefalseAnswer(questionID int, expectAnswer bool) (*models.Truefalse, error)
	CreateShortanswerAnswer(questionID int, expectAnswer string) (*models.Shortanswer, error)
	GetAllAnswers() ([]*models.Truefalse, []*models.Shortanswer, error)
	EditTruefalseAnswer(id int, expectAnswer bool) (*models.Truefalse, error)
	EditShortanswerAnswer(id int, expectAnswer string) (*models.Shortanswer, error)
	DeleteAnswer(id int) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{DB: db}
}

func (r *repository) CreateTruefalseAnswer(questionID int, expectAnswer bool) (*models.Truefalse, error) {
	answer := &models.Truefalse{
		QuestionID:   questionID,
		ExpectAnswer: expectAnswer,
	}
	if err := r.DB.Create(answer).Error; err != nil {
		return nil, err
	}
	return answer, nil
}

func (r *repository) CreateShortanswerAnswer(questionID int, expectAnswer string) (*models.Shortanswer, error) {
	answer := &models.Shortanswer{
		QuestionID:   questionID,
		ExpectAnswer: expectAnswer,
	}
	if err := r.DB.Create(answer).Error; err != nil {
		return nil, err
	}
	return answer, nil
}

func (r *repository) GetAllAnswers() ([]*models.Truefalse, []*models.Shortanswer, error) {
	var truefalseAnswers []*models.Truefalse
	var shortanswerAnswers []*models.Shortanswer

	if err := r.DB.Find(&truefalseAnswers).Error; err != nil {
		return nil, nil, err
	}

	if err := r.DB.Find(&shortanswerAnswers).Error; err != nil {
		return nil, nil, err
	}

	return truefalseAnswers, shortanswerAnswers, nil
}

func (r *repository) EditTruefalseAnswer(id int, expectAnswer bool) (*models.Truefalse, error) {
	var answer models.Truefalse
	if err := r.DB.First(&answer, id).Error; err != nil {
		return nil, err
	}

	answer.ExpectAnswer = expectAnswer
	if err := r.DB.Save(&answer).Error; err != nil {
		return nil, err
	}

	return &answer, nil
}

func (r *repository) EditShortanswerAnswer(id int, expectAnswer string) (*models.Shortanswer, error) {
	var answer models.Shortanswer
	if err := r.DB.First(&answer, id).Error; err != nil {
		return nil, err
	}

	answer.ExpectAnswer = expectAnswer
	if err := r.DB.Save(&answer).Error; err != nil {
		return nil, err
	}

	return &answer, nil
}

func (r *repository) DeleteAnswer(id int) error {
	var answer models.Truefalse
	if err := r.DB.First(&answer, id).Error; err != nil {
		var shortanswer models.Shortanswer
		if err := r.DB.First(&shortanswer, id).Error; err != nil {
			return err
		}
		r.DB.Delete(&shortanswer)
	} else {
		r.DB.Delete(&answer)
	}
	return nil
}

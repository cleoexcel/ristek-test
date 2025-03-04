package submission

import (
	"github.com/cleoexcel/ristek-test/app/models"
	"gorm.io/gorm"
	"errors"
)

type SubmissionRepository interface {
	CreateSubmission(tryoutID int, userID int) (*models.Submission, error)
	GetSubmissionByTryoutID(tryoutID int) (*models.Submission, error)
}

type submissionRepository struct {
	DB *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{DB: db}
}

func (r *submissionRepository) CreateSubmission(tryoutID int, userID int) (*models.Submission, error) {
	var submission models.Submission

	err := r.DB.Where("tryout_id = ?", tryoutID).First(&submission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			submission = models.Submission{
				TryoutID:        tryoutID,
				UserID:          userID,
				NumberOfAttempt: 1, 
			}
			if err := r.DB.Create(&submission).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		submission.NumberOfAttempt++
		if err := r.DB.Save(&submission).Error; err != nil {
			return nil, err
		}
	}

	if err := r.DB.Preload("Tryout").First(&submission, submission.ID).Error; err != nil {
		return nil, err
	}

	return &submission, nil
}


func (r *submissionRepository) GetSubmissionByTryoutID(tryoutID int) (*models.Submission, error) {
	var submission models.Submission
	err := r.DB.Preload("Tryout").Where("tryout_id = ?", tryoutID).First(&submission).Error
	if err != nil {
		return nil, err
	}
	return &submission, nil
}

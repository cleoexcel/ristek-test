package submission

import (
	"github.com/cleoexcel/ristek-test/app/models"
	"gorm.io/gorm"
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
	submission := &models.Submission{
		TryoutID: tryoutID,
		UserID:   userID,
	}

	if err := r.DB.Create(submission).Error; err != nil {
		return nil, err
	}
	return submission, nil
}

func (r *submissionRepository) GetSubmissionByTryoutID(tryoutID int) (*models.Submission, error) {
	var submission models.Submission
	err := r.DB.Where("tryout_id = ?", tryoutID).First(&submission).Error
	if err != nil {
		return nil, err
	}
	return &submission, nil
}

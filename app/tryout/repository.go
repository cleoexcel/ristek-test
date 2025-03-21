package tryout

import (
	"fmt"
	"time"

	"github.com/cleoexcel/ristek-test/app/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateTryout(title, description string, userId int, category string) (*models.Tryout, error)
	GetAllTryout(title string, category string, createdAt string, userId int) ([]*models.Tryout, error)
	GetDetailTryoutByTryoutID(id int) (*models.Tryout, error)
	EditTryoutByTryoutID(id int, title, description string, userId int) (*models.Tryout, error)
	DeleteTryoutByTryoutID(id int) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{DB: db}
}

func (r *repository) CreateTryout(title, description string, userId int, category string) (*models.Tryout, error) {
	tryout := &models.Tryout{
		Title:       title,
		Description: description,
		UserID:      userId,
		Category:    category,
	}
	if err := r.DB.Create(tryout).Error; err != nil {
		return nil, err
	}
	return tryout, nil
}

func (r *repository) GetAllTryout(title string, category string, createdAt string, userId int) ([]*models.Tryout, error) {
	var tryouts []*models.Tryout
	query := r.DB.Model(&models.Tryout{})

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	if category != "" {
		query = query.Where("category LIKE ?", "%"+category+"%")
	}

	if createdAt != "" {
		query = query.Where("DATE(created_at) = ?", createdAt)
	}

	if userId > 0 {
		query = query.Where("user_id = ?", userId)
	}

	err := query.Find(&tryouts).Error
	if err != nil {
		return nil, err
	}
	return tryouts, nil
}

func (r *repository) GetDetailTryoutByTryoutID(id int) (*models.Tryout, error) {
	var tryout models.Tryout
	err := r.DB.First(&tryout, id).Error
	if err != nil {
		return nil, err
	}
	return &tryout, nil
}

func (r *repository) EditTryoutByTryoutID(id int, title, description string, userId int) (*models.Tryout, error) {
	var tryout models.Tryout
	err := r.DB.First(&tryout, id).Error
	if err != nil {
		return nil, err
	}

	if tryout.UserID != userId {
		return nil, fmt.Errorf("you are not authorized to edit this tryout")
	}

	tryout.Title = title
	tryout.Description = description
	now := time.Now()
	tryout.UpdatedAt = &now
	if err := r.DB.Save(&tryout).Error; err != nil {
		return nil, err
	}

	return &tryout, nil
}

func (r *repository) DeleteTryoutByTryoutID(id int) error {
	var tryout models.Tryout
	err := r.DB.First(&tryout, id).Error
	if err != nil {
		return err
	}

	err = r.DB.Delete(&tryout).Error
	if err != nil {
		return err
	}
	return nil
}

package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"github.com/cleoexcel/ristek-test/app/models"
)

type Repository interface {
	GetAllUsers() ([]*models.User, error)
	CreateUser(username, password string) (*models.User, error)
	LoginUser(username, password string) error
	GetUserByUsername(username string) (*models.User, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{DB: db}
}

func (r *repository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil 
		}
		return nil, err
	}

	
	return &user, nil
}

func (r *repository) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	err := r.DB.Find(&users).Error 
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) CreateUser(username, password string) (*models.User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	if err := r.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *repository) LoginUser(username, password string) error {
	var user models.User

	
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid username or password")
		}
		return fmt.Errorf("failed to query user: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return fmt.Errorf("invalid username or password")
	}

	return nil
}

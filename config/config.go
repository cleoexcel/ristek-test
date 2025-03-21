package config

import (
	"fmt"
	"log"
	"os"

	"strconv"

	"github.com/cleoexcel/ristek-test/app/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	JWT_SECRET_KEY    string
	JWT_EXPIRY_IN_DAY int
)

func InitDatabase() (db *gorm.DB) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_TIMEZONE"))

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Tryout{})
	db.AutoMigrate(&models.Question{})
	db.AutoMigrate(&models.Submission{})
	db.AutoMigrate(&models.TrueFalse{})
	db.AutoMigrate(&models.ShortAnswer{})
	db.AutoMigrate(&models.MultipleChoice{})
	db.AutoMigrate(&models.MultipleChoiceOption{})
	db.AutoMigrate(&models.SubmissionAnswerTrueFalse{})
	db.AutoMigrate(&models.SubmissionAnswerShortAnswer{})
	db.AutoMigrate(&models.SubmissionAnswerMultipleChoice{})

	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	JWT_EXPIRY_IN_DAY, _ = strconv.Atoi(os.Getenv("JWT_EXPIRY_IN_DAYS"))

	return
}

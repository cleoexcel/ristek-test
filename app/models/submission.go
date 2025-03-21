package models

import (
	"time"
)

type Submission struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	TryoutID        int       `json:"tryout_id"`
	UserID          int       `json:"user_id"`
	NumberOfAttempt int       `json:"number_of_attempt"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	TotalScore      float64   `json:"total_score"`

	Tryout Tryout `json:"tryout" gorm:"foreignKey:TryoutID"`
}

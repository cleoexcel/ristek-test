package models

import (
	"time"
)

type Submission struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	TryoutID  int       `json:"tryout_id"`   
	UserID    int       `json:"user_id"`     
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Tryout    Tryout    `json:"tryout"`      
	User      User      `json:"user"`         
}

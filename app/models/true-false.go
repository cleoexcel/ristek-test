package models

type TrueFalse struct {
	ID           int  `json:"id" gorm:"primaryKey;autoIncrement"`
	QuestionID   int  `json:"question_id"`
	ExpectAnswer bool `json:"expectanswer"`
}

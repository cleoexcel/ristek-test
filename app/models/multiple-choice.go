package models

type MultipleChoice struct {
	ID          int      `json:"id" gorm:"primaryKey;autoIncrement"`
	QuestionID  int      `json:"question_id"`
	Options     []string `json:"options" gorm:"type:json"`
	ExpectAnswer string `json:"expectanswer"`
}

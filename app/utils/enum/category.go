package enum

func IsValidCategory(category string) bool {
    switch category {
    case
        "Biology",
        "Physics",
        "Math",
        "Chemistry",
		"History":
            return true
    }
    return false
}

func IsValidQuestionType(questionType string) bool {
    switch questionType {
    case
        "T/F",
        "Multiple Choice",
        "Short Answer":
            return true
    }
    return false
}
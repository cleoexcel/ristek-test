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

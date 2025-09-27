package repository

// stringPtr returns a pointer to the given string
func stringPtr(s string) *string {
	return &s
}

// intPtr returns a pointer to the given int
func intPtr(i int) *int {
	return &i
}

// boolPtr returns a pointer to the given bool
func boolPtr(b bool) *bool {
	return &b
}
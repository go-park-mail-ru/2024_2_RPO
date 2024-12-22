package misc

// StringPtr превращает строку в указатель
func StringPtr(s string) *string {
	return &s
}

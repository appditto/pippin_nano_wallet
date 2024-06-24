package utils

// Convert value to pointer with generic interface
func ToPtr[T interface{}](value T) *T {
	return &value
}

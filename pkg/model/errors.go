package model

// UserNotFoundError represents when an instance is not found.
type UserNotFoundError struct{}

func (e UserNotFoundError) Error() string {
	return "User not found"
}

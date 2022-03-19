package api

// InternalServerError represents when an instance is not found.
type InternalServerError struct{}

func (e InternalServerError) Error() string {
	return "Internal Server Error"
}

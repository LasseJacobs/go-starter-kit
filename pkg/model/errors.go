package model

// IsNotFoundError returns whether an error represents a "not found" error.
func IsNotFoundError(err error) bool {
	switch err.(type) {
	case StoryNotFoundError:
		return true
	}
	return false
}

// StoryNotFoundError represents when an instance is not found.
type StoryNotFoundError struct{}

func (e StoryNotFoundError) Error() string {
	return "User not found"
}

// ValidationError used for when a struct does not ahere to all invariants.
type ValidationError struct {
	Violation string
}

func (e ValidationError) Error() string {
	return "User not found"
}

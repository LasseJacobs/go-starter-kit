package model

import (
	"time"
)

// Story represents an entry in our system.
type Story struct {
	ID              string    `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Author          string    `json:"author" db:"author"`
	Votes           int       `json:"votes" db:"votes"`
	Url             string    `json:"url" db:"url"`
	PublicationDate time.Time `json:"originDate" db:"origin_date"`
}

func (s *Story) Vote() {
	s.Votes++
}

func ValidateStory(story *Story) error {
	//TODO: improve this
	if story.ID == "" {
		return ValidationError{
			Violation: "missing ID attribute",
		}
	}
	return nil
}

package functions

import (
	"database/sql"
	"github.com/LasseJacobs/go-starter-kit/internal/storage"
	"github.com/LasseJacobs/go-starter-kit/pkg/model"
	"github.com/pkg/errors"
)

func FindStoryById(tx storage.Connection, id string) (*model.Story, error) {
	var story model.Story
	err := tx.Get(&story, `SELECT id, title, author, votes, url, origin_date FROM content.stories WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.StoryNotFoundError{}
		}
		return nil, errors.Wrap(err, "failed to find story")
	}
	return &story, err
}

// PageStories returns a list of stories, paginated
func PageStories(tx storage.Connection, p model.Pagination) (*[]model.Story, error) {
	var stories []model.Story
	err := tx.Select(&stories, `SELECT id, title, author, votes, url, origin_date 
		FROM content.stories
		ORDER BY transaction_date DESC
		LIMIT $1 OFFSET $2`, p.PerPage, p.Offset())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.StoryNotFoundError{}
		}
		return nil, errors.Wrap(err, "failed to find story")
	}
	return &stories, err
}

func CreateStory(tx storage.Connection, story *model.Story) error {
	_, err := tx.NamedExec("INSERT INTO content.stories (id, title, author, votes, url, origin_date) "+
		"VALUES (:id, :title, :author, :votes, :url, :origin_date)", story)
	if err != nil {
		return errors.Wrap(err, "failed to create story")
	}
	return nil
}

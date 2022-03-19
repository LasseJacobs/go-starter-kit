package api

import (
	"github.com/LasseJacobs/go-starter-kit/internal/functions"
	"github.com/LasseJacobs/go-starter-kit/pkg/model"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// getStory expects {storyid} parameter, return story matching that id
func (a *API) getStory(w http.ResponseWriter, r *http.Request) {
	var storyId = chi.URLParam(r, "storyid")
	story, err := functions.FindStoryById(a.db, storyId)
	if err != nil {
		sendError(w, err)
		return
	}

	sendJSON(w, http.StatusOK, story)
}

// postStory inserts a story into the database if it does not exist yet returning 201, 200 if existing
func (a *API) postStory(w http.ResponseWriter, r *http.Request) {
	var story = &model.Story{}
	err := readModel(r, story)
	if err != nil {
		sendError(w, err)
		return
	}
	err = model.ValidateStory(story)
	if err != nil {
		sendError(w, err)
		return
	}

	//check if story exists already
	old, err := functions.FindStoryById(a.db, story.ID)
	if err != nil && !model.IsNotFoundError(err) {
		sendError(w, err)
		return
	}
	//story already exists, exiting
	if old != nil {
		sendJSON(w, http.StatusOK, old)
		return
	}

	err = functions.CreateStory(a.db, story)
	if err != nil {
		sendError(w, err)
		return
	}
	sendJSON(w, http.StatusCreated, story)
}

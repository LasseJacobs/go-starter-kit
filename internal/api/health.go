package api

import (
	"fmt"
	"net/http"
)

func (a *API) healthCheck(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, http.StatusOK, map[string]string{
		"version":     a.version,
		"name":        a.config.Name,
		"description": a.config.Desc,
	})
}

// failureCheck example endpoint to test application failure
func (a *API) failureCheck(w http.ResponseWriter, r *http.Request) {
	sendError(w, fmt.Errorf("this is some failure"))
}

// failureCheck example endpoint to test application failure
func (a *API) panicCheck(w http.ResponseWriter, r *http.Request) {
	panic(fmt.Errorf("this is some failure"))
}

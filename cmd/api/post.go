package main

import (
	"errors"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Content string   `json:"content"`
	Title   string   `json:"title"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	// TODO: Change after authentication
	userID := 1

	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	post := &store.Post{
		UserID:  int64(userID),
		Content: payload.Content,
		Title:   payload.Title,
		Tags:    payload.Tags,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, *post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ctx := r.Context()

	post, err := app.store.Posts.Get(ctx, int64(id))

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return

		}
	}

	if err := writeJSON(w, http.StatusOK, *post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

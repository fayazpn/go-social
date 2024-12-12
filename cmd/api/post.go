package main

import (
	"context"
	"errors"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Content string   `json:"content" validate:"required,max=100"`
	Title   string   `json:"title" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	// TODO: Change after authentication
	userID := 1

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
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
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusCreated, *post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post, ok := getPostFromCtx(r)
	if !ok {
		app.notFoundResponse(w, r, errors.New("not found in middleware"))
		return
	}

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.writeResponse(w, http.StatusOK, *post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post, ok := getPostFromCtx(r)
	if !ok {
		app.notFoundResponse(w, r, errors.New("not found in middleware"))
		return
	}

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.conflictResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.writeResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post, ok := getPostFromCtx(r)
	if !ok {
		app.notFoundResponse(w, r, errors.New("not found in middleware"))
		return
	}

	err := app.store.Posts.Delete(r.Context(), post)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			params := chi.URLParam(r, "postID")
			id, err := strconv.ParseInt(params, 10, 64)
			if err != nil {
				app.internalServerError(w, r, err)
				return
			}

			ctx := r.Context()

			post, err := app.store.Posts.Get(ctx, int64(id))

			if err != nil {
				switch {
				case errors.Is(err, store.ErrNotFound):
					app.notFoundResponse(w, r, err)
				default:
					app.internalServerError(w, r, err)
				}

				return
			}

			ctx = context.WithValue(ctx, postCtx, post)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
}

func getPostFromCtx(r *http.Request) (*store.Post, bool) {
	post, ok := r.Context().Value(postCtx).(*store.Post)

	return post, ok
}

package main

import (
	"context"
	"errors"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromCtx(r)
	if !ok {
		app.notFoundResponse(w, r, errors.New("user not found in context"))
		return
	}

	if err := app.writeResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser, ok := getUserFromCtx(r)
	if !ok {
		app.notFoundResponse(w, r, errors.New("user not found in context"))
		return
	}

	ctx := r.Context()

	// TODO: Change after impletementing auth
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err := app.store.Followers.Follow(ctx, followerUser.ID, payload.UserID)
	if err != nil {
		// handle error
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowerUser, ok := getUserFromCtx(r)
	if !ok {
		app.notFoundResponse(w, r, errors.New("user not found in context"))
		return
	}

	ctx := r.Context()

	// TODO: Change after impletementing auth
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err := app.store.Followers.Unfollow(ctx, unfollowerUser.ID, payload.UserID)
	if err != nil {
		// handle error
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(params, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetUserByID(ctx, id)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) (*store.User, bool) {
	user, ok := r.Context().Value(userCtx).(*store.User)

	return user, ok
}

package main

import (
	"errors"
	"net/http"

	"github.com/ra-khalish/gosocial/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content"`
	UserID  int64  `json:"user_id"`
	Like    int64  `json:"like"`
}

func (app *application) createCommentPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate post
	if err := validateRequired(payload.Content, "content"); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := validateMax(payload.Content, 1000, "content"); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), payload.UserID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	comment := &store.Comment{
		PostID:  post.ID,
		UserID:  payload.UserID,
		Content: payload.Content,
	}
	// post := &store.Post{
	// 	Title:   payload.Title,
	// 	Content: payload.Content,
	// 	Tags:    payload.Tags,
	// 	// TODO: change after auth
	// 	UserID: 1,
	// }

	if err := app.store.Comment.Create(r.Context(), comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	comment.User = *user

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

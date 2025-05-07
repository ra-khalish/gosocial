package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ra-khalish/gosocial/internal/store"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (p CreatePostPayload) validate() error {
	if err := validateRequired(p.Title, "title"); err != nil {
		return err
	}
	if err := validateRequired(p.Content, "content"); err != nil {
		return err
	}

	if err := validateMin(p.Title, 5, "title"); err != nil {
		return err
	}

	if err := validateMax(p.Title, 100, "content"); err != nil {
		return err
	}

	if err := validateMax(p.Content, 1000, "content"); err != nil {
		return err
	}
	if err := validateMin(p.Content, 5, "content"); err != nil {
		return err
	}
	return nil
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := payload.validate(); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: change after auth
		UserID: 1,
	}

	ctx := r.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	// idParam := chi.URLParam(r, "postID")
	// id, err := strconv.ParseInt(idParam, 10, 64)
	// if err != nil {
	// 	app.internalServerError(w, r, err)
	// 	return
	// }
	// ctx := r.Context()

	// post, err := app.store.Posts.GetByID(ctx, id)
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, store.ErrNotFound):
	// 		app.notFoundResponse(w, r, err)
	// 	default:
	// 		app.internalServerError(w, r, err)
	// 	}
	// 	return
	// }
	post := getPostFromCtx(r)

	comments, err := app.store.Comment.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	// idParam := chi.URLParam(r, "postID")
	// id, err := strconv.ParseInt(idParam, 10, 64)
	// if err != nil {
	// 	app.internalServerError(w, r, err)
	// 	return
	// }
	// ctx := r.Context()
	post := getPostFromCtx(r)

	rows, err := app.store.Posts.Delete(r.Context(), post.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	res := Response{
		TotalData:   rows,
		SuccessData: rows,
		FailData:    0,
	}

	if err := app.jsonResponse(w, http.StatusAccepted, res); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	// w.WriteHeader(http.StatusNoContent)

}

type UpdatePostPayload struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate patch
	validation := func(UpdatePostPayload) error {
		if payload.Content != nil {
			err := validateMax(*payload.Content, 1000, "content")
			if err != nil {
				return err
			}
			post.Content = *payload.Content
		}

		if payload.Title != nil {
			err := validateMin(*payload.Title, 5, "title")
			if err != nil {
				return err
			}
			err = validateMax(*payload.Title, 100, "title")
			if err != nil {
				return err
			}
			post.Title = *payload.Title

		}
		return nil
	}

	if err := validation(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate patch

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := app.store.Posts.GetByID(ctx, id)
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

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}

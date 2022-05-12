package handler

import (
	"encoding/json"
	"net/http"
	"user-service/internal/model"

	"github.com/rs/zerolog"
)

const (
	PostsPath      = "/posts"
	PostPath       = "/post"
	DeletePostPath = "/deletepost"
)

type Posts struct {
	logger  *zerolog.Logger
	service PostsService
}

func NewPosts(logger *zerolog.Logger, srv PostsService) *Posts {
	return &Posts{
		logger:  logger,
		service: srv,
	}
}

type PostsService interface {
	GetPostsFeedByID(string) ([]model.Post, error)
	SubmitPost(string, model.PostSubmitRequest) error
	DeletePost(string, string) error
}

func (h *Posts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid incoming data")
		writeResponse(w, http.StatusUnauthorized, model.Error{Error: "Unauthorized"})
		return
	}
	if r.Method == http.MethodGet {
		posts, err := h.service.GetPostsFeedByID(userID)
		if err != nil {
			h.logger.Error().Err(err).Msg("Unable to get posts")
			writeResponse(w, http.StatusInternalServerError, model.Error{Error: "Internal server error"})
			return
		}
		writeResponse(w, http.StatusOK, posts)
	}
	if r.Method == http.MethodPost {
		if r.URL.String() == PostPath {
			var post model.PostSubmitRequest
			err = json.NewDecoder(r.Body).Decode(&post)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			err := h.service.SubmitPost(userID, post)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusInternalServerError, model.Error{Error: "Error while submitting post"})
				return
			}
		}
		if r.URL.String() == DeletePostPath {
			var postToRemove model.PostDeleteRequest
			err = json.NewDecoder(r.Body).Decode(&postToRemove)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
				return
			}
			err := h.service.DeletePost(postToRemove.PostID.String(), userID)
			if err != nil {
				h.logger.Error().Err(err).Msg("Invalid incoming data")
				writeResponse(w, http.StatusInternalServerError, model.Error{Error: "Error while deleting post"})
				return
			}
		}

	}
}

package handler

import (
	"net/http"
	"strings"
	"user-service/internal/model"

	"github.com/rs/zerolog"
)

const (
	ContentPath = "/content"
)

type Content struct {
	logger  *zerolog.Logger
	service ContentService
}

func NewContent(logger *zerolog.Logger, srv ContentService) *Content {
	return &Content{
		logger:  logger,
		service: srv,
	}
}

type Response struct {
	Content string `json:"content"`
}

type ContentService interface {
	GetUserFeed(string) ([]model.Post, error)
}

func (h *Content) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token := strings.SplitN(authHeader, " ", 2)
	posts, err := h.service.GetUserFeed(token[1])
	if err != nil {

	}
	writeResponse(w, http.StatusOK, posts)
}

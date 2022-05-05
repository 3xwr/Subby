package handler

import (
	"net/http"

	"github.com/rs/zerolog"
)

const (
	ResourcePath = "/test"
)

type Resource struct {
	logger *zerolog.Logger
}

func NewResource(logger *zerolog.Logger) *Resource {
	return &Resource{
		logger: logger,
	}
}

type Response struct {
	Content string `json:"content"`
}

func (h *Resource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Content: "haha you found me",
	}
	writeResponse(w, http.StatusOK, resp)
}

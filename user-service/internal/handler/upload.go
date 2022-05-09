package handler

import (
	"mime/multipart"
	"net/http"
	"user-service/internal/model"

	"github.com/rs/zerolog"
)

const (
	UploadPath = "/upload"
)

type Upload struct {
	logger  *zerolog.Logger
	service UploadService
}

func NewUpload(logger *zerolog.Logger, srv UploadService) *Upload {
	return &Upload{
		logger:  logger,
		service: srv,
	}
}

type UploadService interface {
	UploadImage(file multipart.File, handler *multipart.FileHeader) (string, error)
}

func (h *Upload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("img")
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid incoming data")
		writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
		return
	}
	defer file.Close()

	name, err := h.service.UploadImage(file, handler)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid incoming data")
		writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad Request"})
		return
	}

	resp := model.UploadResponse{
		FileAddress: name,
		Uploaded:    true,
	}

	writeResponse(w, http.StatusOK, resp)
}

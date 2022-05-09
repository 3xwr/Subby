package service

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Upload struct {
	logger *zerolog.Logger
	repo   UploadRepo
}

func NewUpload(logger *zerolog.Logger, repo UploadRepo) *Upload {
	return &Upload{
		logger: logger,
		repo:   repo,
	}
}

const (
	imgBasePath = "../../Subby-images/"
)

var allowedExtensions = []string{".png", ".jpg", ".jpeg"}

type UploadRepo interface {
}

func (s *Upload) UploadImage(file multipart.File, handler *multipart.FileHeader) (string, error) {
	extensionAllowed := false
	fileName := imgBasePath + handler.Filename
	ext := filepath.Ext(fileName)
	for _, extension := range allowedExtensions {
		if ext == extension {
			extensionAllowed = true
		}
	}

	if !extensionAllowed {
		return "", fmt.Errorf("unsupported file format uploaded")
	}

	rand.Seed(time.Now().UnixNano())
	for fileExists(fileName) {
		if len(fileName) > 512 {
			return "", fmt.Errorf("filename is too long")
		}
		fileName = strings.TrimSuffix(fileName, ext)
		fileName += fmt.Sprint(rand.Intn(10))
		fileName += ext
	}

	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	contentTypeAllowed := false

	contentType := http.DetectContentType(fileBytes)
	for _, extension := range allowedExtensions {
		imgType := "image/" + trimLeftChar(extension)
		if contentType == imgType {
			contentTypeAllowed = true
		}
	}

	if contentTypeAllowed {
		f.Write(fileBytes)
	} else {
		return "", fmt.Errorf("unsupported file format uploaded")
	}

	return fileName, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

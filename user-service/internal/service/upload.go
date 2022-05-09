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

type UploadRepo interface {
}

func (s *Upload) UploadImage(file multipart.File, handler *multipart.FileHeader) (string, error) {
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	fileName := imgBasePath + handler.Filename

	if !fileExists(fileName) {
		f, err := os.Create(fileName)
		if err != nil {
			return "", err
		}
		defer f.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return "", err
		}

		contentType := http.DetectContentType(fileBytes)
		if contentType == "image/png" {
			f.Write(fileBytes)
		} else {
			return "", err
		}

		fmt.Println("FILE UPLOADED")
		return fileName, nil
	} else {
		rand.Seed(time.Now().UnixNano())
		name := fileName
		for fileExists(name) {
			if len(name) > 512 {
				return "", fmt.Errorf("filename is too long")
			}
			ext := filepath.Ext(name)
			name = strings.TrimSuffix(name, ext)
			name += fmt.Sprint(rand.Intn(10))
			name += ext
		}
		f, err := os.Create(name)
		if err != nil {
			return "", err
		}
		defer f.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return "", err
		}

		contentType := http.DetectContentType(fileBytes)
		if contentType == "image/png" {
			f.Write(fileBytes)
		} else {
			return "", err
		}

		fmt.Println("FILE UPLOADED")
		return name, nil
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

package services

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/suryansh74/zomato/services/utils-service/internal/adapters"
)

const MaxUploadSize = 2 << 20 // 2MB

type UtilsService struct {
	storage adapters.StorageAdapter // Uses the interface!
}

// NewUtilsService injects the chosen storage adapter
func NewUtilsService(storage adapters.StorageAdapter) *UtilsService {
	return &UtilsService{
		storage: storage,
	}
}

func (s *UtilsService) ProcessAndUploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	// 1. Validate Size
	if header.Size > MaxUploadSize {
		return "", errors.New("image file size must be 2MB or less")
	}

	// 2. Validate Type
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil && err != io.EOF {
		return "", errors.New("failed to read file content")
	}

	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" {
		return "", errors.New("invalid file type. Only JPG/JPEG and PNG are allowed")
	}

	// 3. Reset pointer
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", errors.New("failed to process file stream")
	}

	// 4. Delegate to the injected adapter
	return s.storage.UploadImage(ctx, file, header.Filename)
}

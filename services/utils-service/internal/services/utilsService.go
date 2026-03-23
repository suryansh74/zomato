package services

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// MaxUploadSize is 2 Megabytes
const MaxUploadSize = 2 << 20

type UtilsService struct{}

func NewUtilsService() *UtilsService {
	return &UtilsService{}
}

// ProcessAndUploadImage validates the file size and type, then uploads it to Cloudinary
func (s *UtilsService) ProcessAndUploadImage(ctx context.Context, cld *cloudinary.Cloudinary, file multipart.File, header *multipart.FileHeader) (string, error) {
	// 1. Validate File Size
	if header.Size > MaxUploadSize {
		return "", errors.New("image file size must be 2MB or less")
	}

	// 2. Validate File Type securely (read first 512 bytes)
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil && err != io.EOF {
		return "", errors.New("failed to read file content")
	}

	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" {
		return "", errors.New("invalid file type. Only JPG/JPEG and PNG are allowed")
	}

	// 3. Reset the file pointer back to the beginning
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", errors.New("failed to process file stream")
	}

	// 4. Upload to Cloudinary
	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         "zomato_clone/uploads",
		UniqueFilename: api.Bool(true),
		Overwrite:      api.Bool(false),
	})
	if err != nil {
		return "", errors.New("failed to upload image to cloud storage")
	}

	return resp.SecureURL, nil
}

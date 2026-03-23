package adapters

import (
	"context"
	"io"
)

// StorageAdapter defines the contract for any cloud storage provider
type StorageAdapter interface {
	UploadImage(ctx context.Context, file io.Reader, filename string) (string, error)
}

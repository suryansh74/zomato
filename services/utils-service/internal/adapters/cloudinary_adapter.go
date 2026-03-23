package adapters

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryAdapter struct {
	client *cloudinary.Cloudinary
}

func NewCloudinaryAdapter(cld *cloudinary.Cloudinary) *CloudinaryAdapter {
	return &CloudinaryAdapter{client: cld}
}

// UploadImage implements the StorageAdapter interface
func (c *CloudinaryAdapter) UploadImage(ctx context.Context, file io.Reader, filename string) (string, error) {
	resp, err := c.client.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         "zomato_clone/uploads",
		PublicID:       filename, // Optional: use the original filename
		UniqueFilename: api.Bool(true),
		Overwrite:      api.Bool(false),
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

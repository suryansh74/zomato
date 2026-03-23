package client

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
)

func NewCloudinary(cloudinaryURL string) (*cloudinary.Cloudinary, context.Context, error) {
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, nil, err
	}

	cld.Config.URL.Secure = true
	ctx := context.Background()

	return cld, ctx, nil
}

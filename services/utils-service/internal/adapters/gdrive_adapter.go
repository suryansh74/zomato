package adapters

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDriveAdapter struct {
	service  *drive.Service
	folderID string
}

func NewGoogleDriveAdapter(ctx context.Context, credentialsFile string, folderID string) (*GoogleDriveAdapter, error) {
	// ADD drive.DriveScope HERE
	srv, err := drive.NewService(ctx, option.WithCredentialsFile(credentialsFile), option.WithScopes(drive.DriveScope))
	if err != nil {
		return nil, err
	}
	return &GoogleDriveAdapter{service: srv, folderID: folderID}, nil
}

func (g *GoogleDriveAdapter) UploadImage(ctx context.Context, file io.Reader, filename string) (string, error) {
	fmt.Printf("DEBUG: Attempting upload with Folder ID: '%s'\n", g.folderID)

	f := &drive.File{
		Name: filename,
		// TEMPORARILY COMMENT OUT THE PARENTS LINE TO TEST ROOT UPLOAD
		Parents: []string{g.folderID},
	}

	// Upload the file to Drive
	res, err := g.service.Files.Create(f).
		Media(file).
		SupportsAllDrives(true). // Keep this!
		Context(ctx).
		Do()
	if err != nil {
		fmt.Printf("DEBUG: Upload failed with error: %v\n", err)
		return "", err
	}

	fmt.Printf("DEBUG: Upload succeeded! File ID: %s\n", res.Id)

	permission := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}

	_, err = g.service.Permissions.Create(res.Id, permission).
		SupportsAllDrives(true).
		Context(ctx).
		Do()
	if err != nil {
		fmt.Printf("DEBUG: Permission update failed: %v\n", err)
		return "", err
	}

	viewURL := fmt.Sprintf("https://drive.google.com/uc?id=%s", res.Id)
	return viewURL, nil
}

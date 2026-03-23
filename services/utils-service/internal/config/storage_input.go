package config

type StorageProvider string

const (
	Cloudinary StorageProvider = "cloudinary"
	GDrive     StorageProvider = "gdrive"
)

func (s StorageProvider) IsValid() bool {
	switch s {
	case Cloudinary, GDrive:
		return true
	default:
		return false
	}
}

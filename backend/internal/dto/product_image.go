package dto

import (
	"fmt"
	"io"
)

// FileUpload represents an uploaded file with its metadata.
type FileUpload struct {
	Content     io.Reader // File content stream
	Name        string    // Original filename
	Size        int64     // File size in bytes
	ContentType string    // MIME type (image/jpeg", "image/png")
}

// Validate checks if the file upload meets the application requirements.
func (f *FileUpload) Validate() error {
	// Check file size constraint
	if f.Size > 5*1024*1024 {
		return fmt.Errorf("image too large (max 5MB)")
	}

	// Validate against allowed image formats
	allowedTypes := []string{"image/jpeg", "image/png"}
	for _, allowedType := range allowedTypes {
		if f.ContentType == allowedType {
			return nil
		}
	}

	return fmt.Errorf("invalid image type: %s", f.ContentType)
}

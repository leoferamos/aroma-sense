package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// deriveThumbnailKey creates a thumbnail key from the original name: foo/bar.jpg -> foo/bar_thumb.jpg
func deriveThumbnailKey(name string) string {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	return base + "_thumb.jpg"
}

// resizeToFitJPEG performs a simple proportional resize to fit within (maxW,maxH) and encodes as JPEG.
func resizeToFitJPEG(r io.Reader, maxW, maxH int, quality int) ([]byte, error) {
	if maxW <= 0 {
		maxW = 256
	}
	if maxH <= 0 {
		maxH = 256
	}
	if quality <= 0 || quality > 100 {
		quality = 80
	}

	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()
	// compute new size keeping aspect ratio
	scale := 1.0
	if w > maxW {
		scale = float64(maxW) / float64(w)
	}
	if int(float64(h)*scale) > maxH {
		scale = float64(maxH) / float64(h)
	}
	newW := int(float64(w) * scale)
	newH := int(float64(h) * scale)
	if newW <= 0 {
		newW = 1
	}
	if newH <= 0 {
		newH = 1
	}

	// simple resize using manual nearest-neighbor
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	// Use nearest-neighbor manually for speed
	for y := 0; y < newH; y++ {
		sy := int(float64(y) * float64(h) / float64(newH))
		for x := 0; x < newW; x++ {
			sx := int(float64(x) * float64(w) / float64(newW))
			dst.Set(x, y, img.At(sx, sy))
		}
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: quality}); err != nil {
		return nil, fmt.Errorf("encode jpeg: %w", err)
	}
	return buf.Bytes(), nil
}

// UploadImageWithThumbnail uploads the original image and a thumbnail variant.
func (s *SupabaseS3) UploadImageWithThumbnail(ctx context.Context, imageName string, content io.Reader, size int64, contentType string, maxW, maxH int) (string, string, error) {
	// Read the content once into memory to both upload original and create thumbnail.
	var originalBuf bytes.Buffer
	if _, err := io.Copy(&originalBuf, content); err != nil {
		return "", "", fmt.Errorf("read original: %w", err)
	}

	//Upload original
	origInput := &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(imageName),
		Body:        bytes.NewReader(originalBuf.Bytes()),
		ContentType: aws.String(contentType),
	}
	if _, err := s.Client.PutObject(ctx, origInput); err != nil {
		return "", "", fmt.Errorf("upload original: %w", err)
	}
	origURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", strings.TrimRight(s.PublicURL, "/"), s.Bucket, imageName)

	//Create and upload thumbnail
	thumbKey := deriveThumbnailKey(imageName)
	thumbBytes, err := resizeToFitJPEG(bytes.NewReader(originalBuf.Bytes()), maxW, maxH, 80)
	if err != nil {
		// If resizing fails, return original URL and empty thumbnail URL.
		return origURL, "", nil
	}
	thumbInput := &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(thumbKey),
		Body:        bytes.NewReader(thumbBytes),
		ContentType: aws.String("image/jpeg"),
	}
	if _, err := s.Client.PutObject(ctx, thumbInput); err != nil {
		// If thumbnail upload fails, still return original.
		return origURL, "", nil
	}
	thumbURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", strings.TrimRight(s.PublicURL, "/"), s.Bucket, thumbKey)
	return origURL, thumbURL, nil
}

// DeleteImage deletes an image from storage
func (s *SupabaseS3) DeleteImage(ctx context.Context, imageName string) error {
	// Delete the original image
	origInput := &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(imageName),
	}
	if _, err := s.Client.DeleteObject(ctx, origInput); err != nil {
		return fmt.Errorf("delete original image: %w", err)
	}

	// Delete the thumbnail if it exists
	thumbKey := deriveThumbnailKey(imageName)
	thumbInput := &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(thumbKey),
	}
	// Don't return error if thumbnail doesn't exist
	if _, err := s.Client.DeleteObject(ctx, thumbInput); err != nil {
		// Log error but don't fail
		fmt.Printf("Warning: failed to delete thumbnail %s: %v\n", thumbKey, err)
	}

	return nil
}

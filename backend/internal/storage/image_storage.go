package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ImageStorage defines an interface for image upload
type ImageStorage interface {
	UploadImage(ctx context.Context, imageName string, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

type SupabaseS3 struct {
	Client   *s3.Client
	Bucket   string
	Endpoint string
}

// NewSupabaseS3 initializes and returns a SupabaseS3 instance, or error if config is invalid
func NewSupabaseS3() (*SupabaseS3, error) {
	endpoint := os.Getenv("SUPABASE_S3_ENDPOINT")
	region := os.Getenv("SUPABASE_S3_REGION")
	accessKey := os.Getenv("SUPABASE_S3_ACCESS_KEY")
	secretKey := os.Getenv("SUPABASE_S3_SECRET_KEY")
	bucket := os.Getenv("SUPABASE_BUCKET")

	// Basic validation of required env variables
	if endpoint == "" || region == "" || accessKey == "" || secretKey == "" || bucket == "" {
		return nil, fmt.Errorf("missing required Supabase S3 environment variables")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &SupabaseS3{
		Client:   s3Client,
		Bucket:   bucket,
		Endpoint: endpoint,
	}, nil
}

// UploadImage uploads an image to Supabase S3 and returns its public URL
func (s *SupabaseS3) UploadImage(ctx context.Context, imageName string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(imageName),
		Body:        file,
		ContentType: aws.String(contentType),
	}

	_, err := s.Client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}

	publicURL := fmt.Sprintf("%s/%s/%s", strings.TrimRight(s.Endpoint, "/"), s.Bucket, imageName)
	return publicURL, nil
}

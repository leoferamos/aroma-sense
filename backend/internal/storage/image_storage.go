package storage

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ImageStorage defines an interface for image upload
type ImageStorage interface {
	UploadImageWithThumbnail(ctx context.Context, imageName string, content io.Reader, size int64, contentType string, maxW, maxH int) (string, string, error)
	DeleteImage(ctx context.Context, imageName string) error
}

type SupabaseS3 struct {
	Client    *s3.Client
	Bucket    string
	PublicURL string
}

// NewSupabaseS3 initializes and returns a SupabaseS3 instance, or error if config is invalid
func NewSupabaseS3() (*SupabaseS3, error) {
	endpoint := os.Getenv("SUPABASE_S3_ENDPOINT")
	region := os.Getenv("SUPABASE_S3_REGION")
	accessKey := os.Getenv("SUPABASE_S3_ACCESS_KEY")
	secretKey := os.Getenv("SUPABASE_S3_SECRET_KEY")
	bucket := os.Getenv("SUPABASE_BUCKET")
	publicURL := os.Getenv("SUPABASE_PUBLIC_URL")

	// Basic validation of required env variables
	if endpoint == "" || region == "" || accessKey == "" || secretKey == "" || bucket == "" || publicURL == "" {
		return nil, fmt.Errorf("missing required Supabase S3 environment variables, including SUPABASE_PUBLIC_URL")
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
		Client:    s3Client,
		Bucket:    bucket,
		PublicURL: publicURL,
	}, nil
}

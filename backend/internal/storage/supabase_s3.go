package storage

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type SupabaseS3 struct {
	Client   *s3.Client
	Bucket   string
	Endpoint string
}

// NewSupabaseS3 initializes and returns a SupabaseS3 instance
func NewSupabaseS3() *SupabaseS3 {
	endpoint := os.Getenv("SUPABASE_S3_ENDPOINT")
	region := os.Getenv("SUPABASE_S3_REGION")
	accessKey := os.Getenv("SUPABASE_S3_ACCESS_KEY")
	secretKey := os.Getenv("SUPABASE_S3_SECRET_KEY")
	bucket := os.Getenv("SUPABASE_BUCKET")

	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			},
		)),
	)
	return &SupabaseS3{
		Client:   s3.NewFromConfig(cfg),
		Bucket:   bucket,
		Endpoint: endpoint,
	}
}

// UploadFile uploads a file to Supabase S3 and returns its public URL
func (s *SupabaseS3) UploadFile(ctx context.Context, imageName string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(imageName),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", err
	}
	publicURL := s.Endpoint + "/" + s.Bucket + "/" + imageName
	return publicURL, nil
}

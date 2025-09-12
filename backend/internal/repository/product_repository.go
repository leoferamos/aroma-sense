import (
	"context"
	"mime/multipart"
	"os"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func newSupabaseS3Client() *s3.Client {
	endpoint := os.Getenv("SUPABASE_S3_ENDPOINT")
	region := os.Getenv("SUPABASE_S3_REGION")
	accessKey := os.Getenv("SUPABASE_S3_ACCESS_KEY")
	secretKey := os.Getenv("SUPABASE_S3_SECRET_KEY")
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
	return s3.NewFromConfig(cfg)
}

type ProductRepository interface {
	Create(input dto.ProductFormDTO, imageName string, file multipart.File, fileHeader *multipart.FileHeader) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(input dto.ProductFormDTO, imageName string, file multipart.File, fileHeader *multipart.FileHeader) error {
	bucket := os.Getenv("SUPABASE_BUCKET")
	client := newSupabaseS3Client()
	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(imageName),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return err
	}
	imageURL := os.Getenv("SUPABASE_S3_ENDPOINT") + "/" + bucket + "/" + imageName
	// TODO: Montar o model Product com imageURL e salvar no banco
	return nil
}

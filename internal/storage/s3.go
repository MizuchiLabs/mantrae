package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"sort"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/settings"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	Bucket       string
	Region       string
	Endpoint     string
	AccessKey    string
	SecretKey    string
	UsePathStyle bool
}

// S3Storage implements StorageBackend for AWS S3
type S3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(ctx context.Context, sm *settings.SettingsManager) (*S3Storage, error) {
	cfgData, err := getSettings(ctx, sm)
	if err != nil {
		return nil, err
	}

	var opts []func(*config.LoadOptions) error

	// Set region
	if cfgData.Region != "" {
		opts = append(opts, config.WithRegion(cfgData.Region))
	}

	// Set credentials
	if cfgData.AccessKey != "" && cfgData.SecretKey != "" {
		creds := credentials.NewStaticCredentialsProvider(cfgData.AccessKey, cfgData.SecretKey, "")
		opts = append(opts, config.WithCredentialsProvider(creds))
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}

	// Build S3 options directly (the preferred modern way)
	s3Options := s3.Options{
		Region:      awsCfg.Region,
		Credentials: awsCfg.Credentials,
	}

	if cfgData.Endpoint != "" {
		parsedURL, err := url.Parse(cfgData.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("invalid endpoint URL: %w", err)
		}

		s3Options.BaseEndpoint = aws.String(cfgData.Endpoint)
		s3Options.EndpointOptions = s3.EndpointResolverOptions{
			DisableHTTPS: parsedURL.Scheme == "http",
		}
	}
	if cfgData.UsePathStyle {
		s3Options.UsePathStyle = true
	}

	return &S3Storage{
		bucket: cfgData.Bucket,
		client: s3.New(s3Options),
	}, nil
}

func (s *S3Storage) Store(ctx context.Context, name string, data io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(name),
		Body:   data,
	})
	if err != nil {
		return fmt.Errorf("failed to store object in S3: %w", err)
	}
	return nil
}

func (s *S3Storage) Retrieve(ctx context.Context, name string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(name),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve object from S3: %w", err)
	}
	return output.Body, nil
}

// List lists the files in S3
func (s *S3Storage) List(ctx context.Context) ([]StoredFile, error) {
	var files []StoredFile
	var nextToken *string

	for {
		output, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:            aws.String(s.bucket),
			ContinuationToken: nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list objects from S3: %w", err)
		}

		for _, object := range output.Contents {
			// Safely dereference pointers.
			key := ""
			if object.Key != nil {
				key = *object.Key
			}
			var ts time.Time
			if object.LastModified != nil {
				ts = *object.LastModified
			}
			files = append(files, StoredFile{
				Name:      key,
				Timestamp: ts,
				Size:      *object.Size,
			})
		}

		nextToken = output.NextContinuationToken

		if nextToken == nil {
			break
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Timestamp.After(files[j].Timestamp)
	})

	return files, nil
}

func (s *S3Storage) Delete(ctx context.Context, name string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(name),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}
	return nil
}

func getSettings(ctx context.Context, sm *settings.SettingsManager) (*S3Config, error) {
	region, err := sm.Get(ctx, settings.KeyS3Region)
	if err != nil {
		return nil, err
	}

	endpoint, err := sm.Get(ctx, settings.KeyS3Endpoint)
	if err != nil {
		return nil, err
	}

	bucket, err := sm.Get(ctx, settings.KeyS3Bucket)
	if err != nil {
		return nil, err
	}

	accessKey, err := sm.Get(ctx, settings.KeyS3AccessKey)
	if err != nil {
		return nil, err
	}

	secretKey, err := sm.Get(ctx, settings.KeyS3SecretKey)
	if err != nil {
		return nil, err
	}

	usePathStyle, err := sm.Get(ctx, settings.KeyS3UsePathStyle)
	if err != nil {
		return nil, err
	}

	cfg := S3Config{
		Region:       region.Value.(string),
		Endpoint:     endpoint.Value.(string),
		Bucket:       bucket.Value.(string),
		AccessKey:    accessKey.Value.(string),
		SecretKey:    secretKey.Value.(string),
		UsePathStyle: usePathStyle.Value.(bool),
	}

	return &cfg, nil
}

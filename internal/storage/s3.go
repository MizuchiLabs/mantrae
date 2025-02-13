package storage

import (
	"context"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Storage implements StorageBackend for AWS S3
type S3Storage struct {
	client *s3.Client
	bucket string
}

// NewS3Storage creates a new S3Storage instance
func NewS3Storage(
	ctx context.Context,
	bucket string,
	awsConfig aws.Config,
) (*S3Storage, error) {
	if bucket == "" {
		return nil, fmt.Errorf("bucket name cannot be empty")
	}

	var cfg aws.Config
	var err error

	if awsConfig.Region == "" {
		cfg, err = config.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS config: %w", err)
		}
	} else {
		cfg = awsConfig
	}

	return &S3Storage{
		bucket: bucket,
		client: s3.NewFromConfig(cfg),
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

package config

import (
	"context"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

type Config struct {
	S3AccessKey        string            `split_words:"true" required:"true"`
	S3Secret           string            `split_words:"true" required:"true"`
	S3Region           string            `split_words:"true" default:"us-east-1"`
	S3Endpoint         string            `split_words:"true"`
	S3Bucket           string            `split_words:"true" required:"true"`
	UploadDir          string            `split_words:"true" default:"/tmp"`
	ListenAddr         string            `split_words:"true" default:"0.0.0.0:8000"`
	Credentials        map[string]string `split_words:"true" required:"true"` // comma separated user:password pairs
	MaxUploadSizeBytes int64             `split_words:"true" default:"10000000"`
}

var (
	once sync.Once
)

// NewS3Client creates new S3 client from config
func (c *Config) NewS3Client(ctx context.Context) (*s3.S3, error) {
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(c.S3AccessKey, c.S3Secret, ""),

		Region:           aws.String(c.S3Region),
		S3ForcePathStyle: aws.Bool(true),
		UseDualStack:     aws.Bool(true),
	}

	if c.S3Endpoint != "" {
		s3Config.Endpoint = aws.String(c.S3Endpoint)

		if strings.HasPrefix(c.S3Endpoint, "http://") {
			s3Config.DisableSSL = aws.Bool(true)
		}
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating AWS session object")
	}
	s3Client := s3.New(newSession)

	// create bucket if it does not exist (on first request)
	var bucketErr error
	once.Do(func() {
		bucketParts := strings.Split(c.S3Bucket, "/")
		var ferr error
		_, ferr = s3Client.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(bucketParts[0]),
		})
		if ferr != nil {
			var s3Err s3.RequestFailure

			if errors.As(ferr, &s3Err) && s3Err.StatusCode() == 404 {
				_, ferr = s3Client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
					Bucket: aws.String(bucketParts[0]),
				})
				if ferr != nil {
					bucketErr = ferr
					return
				}
				return
			}
			bucketErr = ferr
			return
		}
	})
	if bucketErr != nil {
		return nil, errors.Wrap(bucketErr, "Error accessing s3 bucket")
	}

	return s3Client, nil
}

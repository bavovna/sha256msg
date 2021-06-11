package s3store

import (
	"context"
	"io"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-pkgz/repeater"
	"github.com/go-pkgz/repeater/strategy"
	"github.com/pkg/errors"
)

const (
	s3RepeaterFactor = 1.5
	s3RepeatTimes    = 5
)

type s3Client interface {
	PutObjectWithContext(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error)
	GetObjectWithContext(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error)
}

// ErrorNoResults no match for a given key
var ErrorNoResults error = errors.New("Not found")

// upload uploads a file; ReadSeeker is needed by aws lib
func upload(ctx context.Context, s3Client s3Client, bucketWithPath string, blobPath string, r io.ReadSeeker) error {
	pathPrefix := s3prefix(bucketWithPath)

	f := func() error {
		_, err := s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Body:   r,
			Bucket: aws.String(bucket(bucketWithPath)),
			Key:    aws.String(path.Join(pathPrefix, blobPath)),
		})
		if err != nil {
			return errors.Wrap(err, "Error uploading blob to s3")
		}
		return nil
	}

	rp := repeater.New(&strategy.Backoff{
		Repeats: s3RepeatTimes,
		Factor:  s3RepeaterFactor,
		Jitter:  true,
	})
	if err := rp.Do(ctx, f); err != nil {
		return errors.Wrap(err, "Repeater tried hard, but could not upload to S3")
	}
	return nil
}

// download downloads blob from s3
func download(ctx context.Context, s3Client s3Client, bucketWithPath string, blobPath string) (io.ReadCloser, error) {
	pathPrefix := s3prefix(bucketWithPath)

	var res io.ReadCloser
	f := func() error {
		blobRes, err := s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket(bucketWithPath)),
			Key:    aws.String(path.Join(pathPrefix, blobPath)),
		})
		if err != nil {
			return errors.Wrap(err, "Error downloading blob")
		}
		// NOTE: client is responsible of closing
		res = blobRes.Body
		return nil
	}

	rp := repeater.New(&strategy.Backoff{
		Repeats: s3RepeatTimes,
		Factor:  s3RepeaterFactor,
		Jitter:  true,
	})
	if err := rp.Do(ctx, f); err != nil {
		return nil, errors.Wrap(err, "Repeater tried hard, but could not download from S3")
	}
	return res, nil
}

func bucket(bucketWithPath string) string {
	bucketPath := strings.Split(bucketWithPath, "/")
	return bucketPath[0]
}

func s3prefix(bucketWithPath string, pathParts ...string) string {
	bucketPath := strings.Split(bucketWithPath, "/")
	if len(bucketPath) > 1 {
		s3PathParts := append(bucketPath[1:], pathParts...) //nolint: append is OK here
		return path.Join(s3PathParts...)
	}
	return path.Join(pathParts...)
}

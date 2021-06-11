package s3store

import (
	"context"
	"io"

	"github.com/mkorenkov/sha256msg/internal/config"
	"github.com/mkorenkov/sha256msg/internal/sha256sum"
	"github.com/pkg/errors"
)

//API s3 store implementation
type API struct {
	Conf   config.Config
	Client s3Client
}

// Fetch downloads `io.ReadCloser` for the given object key. Client is responsible for calling Close() on the result.
func (s *API) Fetch(ctx context.Context, objectKey string) (io.ReadCloser, error) {
	return download(ctx, s.Client, s.Conf.S3Bucket, objectPath(objectKey))
}

// Store saves `io.ReadSeeker` and returns the object key.
func (s *API) Store(ctx context.Context, r io.ReadSeeker) (string, error) {
	if _, err := r.Seek(0, 0); err != nil {
		return "", errors.Wrap(err, "fseek error")
	}
	objectKey, err := sha256sum.FromReadSeeker(r)
	if err != nil {
		return "", errors.Wrap(err, "hashing error")
	}
	if _, err := r.Seek(0, 0); err != nil {
		return "", errors.Wrap(err, "fseek error")
	}

	return objectKey, upload(ctx, s.Client, s.Conf.S3Bucket, objectPath(objectKey), r)
}

func objectPath(hash string) string {
	if len(hash) < 24 {
		return ""
	}
	// formats key as '68e6/51e67e83/58bef848/everythingelse'
	return hash[0:4] + "/" + hash[4:12] + "/" + hash[12:20] + "/" + hash[20:]
}

package s3store

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	"github.com/mkorenkov/sha256msg/internal/config"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const (
	testAwsKey    = "SECRET-ACCESSKEYID"
	testAwsSecret = "SECRET-SECRETACCESSKEY"
	testAwsRegion = "eu-central-1"
)

func testS3Client(awsEndpoint string) *s3.S3 {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(testAwsKey, testAwsSecret, ""),
		Endpoint:         aws.String(awsEndpoint),
		Region:           aws.String(testAwsRegion),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		panic(errors.Wrap(err, "Error creating AWs session object"))
	}
	return s3.New(newSession)
}

func TestIntegrationS3(t *testing.T) {
	backend := s3mem.New()
	faker := gofakes3.New(backend)
	ts := httptest.NewServer(faker.Server())
	defer ts.Close()

	s3TestBucket := "my-testbucket-upload"

	testCases := []struct {
		content      string
		bucket       string
		expoectedKey string
	}{
		{"Hello World", s3TestBucket, "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e"},
		{"Test Input", fmt.Sprintf("%s/a/b/c", s3TestBucket), "154b35c35faec179a721cbaee303cd9b7460d3da94f4afd1d59c9fa50a8cff87"},
	}

	s3Client := testS3Client(ts.URL)
	_, err := s3Client.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(s3TestBucket)})
	require.NoError(t, err)

	for _, tt := range testCases {
		t.Run(tt.content, func(t *testing.T) {
			c := &API{
				Conf:   config.Config{S3Bucket: tt.bucket},
				Client: s3Client,
			}

			// upload to S3
			key, err := c.Store(context.TODO(), strings.NewReader(tt.content))
			require.NoError(t, err)
			require.Equal(t, tt.expoectedKey, key)

			// validating upload on S3
			s3Path := path.Join(s3prefix(tt.bucket), objectPath(key))

			res1, err := s3Client.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(s3TestBucket),
				Key:    aws.String(s3Path),
			})
			require.NoError(t, err)
			require.NotNil(t, res1)
			defer res1.Body.Close()
			s3Bytes, err := ioutil.ReadAll(res1.Body)
			require.NoError(t, err)
			require.Equal(t, tt.content, string(s3Bytes))

			// fetching using the public API
			res2, err := c.Fetch(context.TODO(), key)
			require.NoError(t, err)
			require.NotNil(t, res2)
			defer res2.Close()

			s3Bytes, err = ioutil.ReadAll(res2)
			require.NoError(t, err)
			require.Equal(t, tt.content, string(s3Bytes))
		})
	}
}

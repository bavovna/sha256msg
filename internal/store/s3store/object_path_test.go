package s3store

import (
	"strings"
	"testing"

	"github.com/mkorenkov/sha256msg/internal/sha256sum"
)

func TestObjectPath(t *testing.T) {
	testCases := map[string]string{
		"Hello World": "a591/a6d40bf4/20404a01/1733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e",
		"Test Input":  "154b/35c35fae/c179a721/cbaee303cd9b7460d3da94f4afd1d59c9fa50a8cff87",
	}

	for input, expectedResult := range testCases {
		objectKey, err := sha256sum.FromReader(strings.NewReader(input))
		if err != nil {
			t.Fatal(err)
		}
		actualResult := objectPath(objectKey)
		if actualResult != expectedResult {
			t.Fatalf("not equal %v != %v", actualResult, expectedResult)
		}
	}
}

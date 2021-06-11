package sha256sum

import (
	"strings"
	"testing"
)

func TestSHA256(t *testing.T) {
	testCases := map[string]string{
		"Hello World": "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e",
		"Test Input":  "154b35c35faec179a721cbaee303cd9b7460d3da94f4afd1d59c9fa50a8cff87",
	}

	for input, expectedResult := range testCases {
		objectKey, err := FromReader(strings.NewReader(input))
		if err != nil {
			t.Fatal(err)
		}
		if objectKey != expectedResult {
			t.Fatalf("not equal %v != %v", objectKey, expectedResult)
		}
	}
}

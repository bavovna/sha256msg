package sha256sum

import (
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// FromReadSeeker calculates sha256 from a given io.ReadSeeker.
func FromReadSeeker(r io.ReadSeeker) (res string, err error) {
	if _, err = r.Seek(0, 0); err != nil {
		return "", errors.Wrap(err, "file seek err")
	}

	res, err = FromReader(r)
	if err != nil {
		return "", err
	}

	if _, err = r.Seek(0, 0); err != nil {
		return "", errors.Wrap(err, "file seek err")
	}

	return res, err
}

// FromReader calculates sha256 from a given io.Reader
func FromReader(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", errors.Wrap(err, "sha256 copy err")
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

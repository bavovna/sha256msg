package constants

import "time"

const (
	// DefaultHTTPServerReadTimeout reasonable default for http.Server ReadTimeout
	DefaultHTTPServerReadTimeout = 10 * time.Second
	// DefaultHTTPServerReadTimeout reasonable default for  http.Server WriteTimeout
	DefaultHTTPServerWriteTimeout = 10 * time.Second
)

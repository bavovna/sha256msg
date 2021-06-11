package middlewares

import "net/http"

// CheckExceedTheLimit validates amount of requests against datasource and returns true if request is allowed.
type CheckExceedTheLimit func(clientIP string) bool

// RateLimit throttles incoming HTTP requests per IP. Helps against DDoS and content scans.
// http://stackoverflow.com/questions/20298220/rate-limiting-http-requests-via-http-handlerfunc-middleware
func RateLimit(h http.Handler, validationFunc CheckExceedTheLimit) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		remoteIP := r.Header.Get("REMOTE_ADDR")
		if validationFunc(remoteIP) {
			w.WriteHeader(http.StatusTooManyRequests)
			// it then returns, not passing the request down the chain
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

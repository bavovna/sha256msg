package requestcontext

import (
	"context"
	"net/http"

	"github.com/mkorenkov/sha256msg/internal/config"
)

type ctxKey struct{}

// RequestContext travels with request
type RequestContext struct {
	Config config.Config
}

// GetRequestContext returns the RequestContext stored in the given context.
func GetRequestContext(ctx context.Context) interface{} {
	return ctx.Value(ctxKey{})
}

// WithContext create a new context from the given one and stores RequestContext object in it.
func WithContext(ctx context.Context, rc interface{}) context.Context {
	return context.WithValue(ctx, ctxKey{}, rc)
}

// InjectRequestContextMiddleware injects a given request context into HTTP request.
func InjectRequestContextMiddleware(handler http.Handler, rc interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := WithContext(r.Context(), rc)
		newReq := r.WithContext(ctx)
		handler.ServeHTTP(w, newReq)
	})
}

// FromRequest returns RequestContext from a given http.Request
func FromRequest(r *http.Request) *RequestContext {
	if val, ok := GetRequestContext(r.Context()).(RequestContext); ok {
		return &val
	}
	return nil
}

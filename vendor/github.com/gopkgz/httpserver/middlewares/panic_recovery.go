package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

// PanicRecoveryMiddleware recovers from panic().
func PanicRecoveryMiddleware(handler http.Handler, reportFunc func(error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if smth := recover(); smth != nil {
				var err error
				if asErr, ok := err.(error); ok {
					err = asErr
				} else {
					err = fmt.Errorf("Non-error panic %v", smth)
				}
				w.WriteHeader(http.StatusInternalServerError)
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, true)
				fmt.Fprintf(w, "INTERNAL SERVER ERROR \n%s\n%s", err, buf)
				log.Printf("[ERROR] %s\n%s", err, buf)
				reportFunc(err)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gopkgz/httpserver/middlewares"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/mkorenkov/sha256msg/internal/config"
	"github.com/mkorenkov/sha256msg/internal/http/noop"
	"github.com/mkorenkov/sha256msg/internal/message"
	"github.com/mkorenkov/sha256msg/internal/requestcontext"
)

func report(err error) {
	log.Printf("[PANIC] %+v\n", err)
}

func run() error {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		return errors.Wrap(err, "error loading config")
	}

	rctx := requestcontext.RequestContext{
		Config: cfg,
	}

	b := middlewares.NewBasicAuthMiddleware(cfg.Credentials)

	r := mux.NewRouter()
	r.HandleFunc("/", noop.Handler)

	api := r.PathPrefix("/api/v1/").Subrouter()
	api.HandleFunc("/", noop.Handler).Methods(http.MethodGet)
	api.HandleFunc("/{key}", message.FetchHandler).Methods(http.MethodGet)
	api.HandleFunc("/", message.StoreHandler).Methods(http.MethodPost)
	api.Use(b.BasicAuth)

	log.Printf("[INFO] Listening %s\n", cfg.ListenAddr)

	srv := &http.Server{
		Handler:           middlewares.LogMiddleware(requestcontext.InjectRequestContextMiddleware(middlewares.PanicRecoveryMiddleware(r, report), rctx)),
		Addr:              cfg.ListenAddr,
		WriteTimeout:      30 * time.Minute,
		ReadTimeout:       30 * time.Minute,
		ReadHeaderTimeout: 10 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		return errors.Wrap(err, "fatal webserver error")
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"embed"
	"io/fs"
	"net/http"
	"time"

	"github.com/mr-karan/doggo/internal/app"
	"github.com/mr-karan/doggo/pkg/utils"
	"github.com/mr-karan/logf"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/knadh/koanf"
)

var (
	// Version of the build. This is injected at build-time.
	buildString = "unknown"

	logger = utils.InitLogger()
	ko     = koanf.New(".")
	//go:embed assets/*
	assetsDir embed.FS
	//go:embed index.html
	html []byte
)

func main() {
	initConfig()

	// Initialize app.
	app := app.New(logger, buildString)

	// Register router instance.
	r := chi.NewRouter()

	// Register middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Frontend Handlers.
	assets, _ := fs.Sub(assetsDir, "assets")
	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/assets/", http.FileServer(http.FS(assets)))
		fs.ServeHTTP(w, r)
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Write(html)
	})

	// API Handlers.
	r.Get("/api/", wrap(app, handleIndexAPI))
	r.Get("/api/ping/", wrap(app, handleHealthCheck))
	r.Post("/api/lookup/", wrap(app, handleLookup))

	// HTTP Server.
	srv := &http.Server{
		Addr:         ko.String("server.address"),
		Handler:      r,
		ReadTimeout:  ko.Duration("server.read_timeout") * time.Millisecond,
		WriteTimeout: ko.Duration("server.write_timeout") * time.Millisecond,
		IdleTimeout:  ko.Duration("server.keepalive_timeout") * time.Millisecond,
	}

	logger.WithFields(logf.Fields{
		"address": srv.Addr,
		"version": buildString,
	}).Info("starting server")

	if err := srv.ListenAndServe(); err != nil {
		logger.WithError(err).Fatal("couldn't start server")
	}
}

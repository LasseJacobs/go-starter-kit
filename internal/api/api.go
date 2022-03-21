package api

import (
	"context"
	"github.com/LasseJacobs/go-starter-kit/internal/config"
	"github.com/LasseJacobs/go-starter-kit/internal/middleware"
	"github.com/LasseJacobs/go-starter-kit/internal/storage"
	"github.com/go-chi/chi/v5"
	chim "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// API is the main REST API
type API struct {
	handler http.Handler
	db      storage.Connection
	config  *config.Config
	version string
}

func NewAPIWithVersion(conf *config.Config, db storage.Connection, version string) *API {
	var api = &API{db: db, config: conf, version: version}
	// setup server routing
	r := chi.NewRouter()

	r.Use(chim.Logger)
	r.Use(chim.Recoverer)
	r.Use(middleware.ResponseLatency)

	// homepage welcome page
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.HTML(w, r, "<html><head><title>Go Starter Kit</title></head><body>Welcome to Go Starter Kit</head></body></html>")
	})

	// register health check route
	r.Get("/health", api.healthCheck)

	// example:
	r.Get("/error", api.failureCheck)
	r.Get("/panic", api.panicCheck)

	// register v1 api path group
	r.Route("/v1", func(r chi.Router) {
		//r.Use(mid.APIVersionCtx("v1"))
		//user.RegisterHandlers(r, db, logger, validate)
		r.Get("/story/{storyid}", api.getStory)
		r.Post("/story", api.postStory)
		r.Get("/story", api.getStories)
	})

	corsHandler := cors.New(cors.Options{
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	api.handler = corsHandler.Handler(r)
	return api
}

func (a *API) Start() error {
	serverErrors := make(chan error, 1)

	server := http.Server{
		Addr:    a.config.Server.Port,
		Handler: a.handler,
	}

	// Start the service listening for requests.
	go func() {
		logrus.Infof("http server listening on :%s", a.config.Server.Port)
		serverErrors <- server.ListenAndServe()
	}()

	// Shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don'usertransport collect this error.
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	// Blocking main and waiting for shutdown.
	case sig := <-shutdown:
		logrus.Infof("start shutdown: %v", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 10)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := server.Shutdown(ctx)
		if err != nil {
			logrus.Infof("graceful shutdown did not complete in %v : %v", 10, err)
			err = server.Close()
			return err
		}

		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}

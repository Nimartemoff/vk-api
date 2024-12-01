package httpserver

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	readHeaderTimeout = 10 * time.Second
	readTimeout       = 600 * time.Second
	idleTimeout       = 180 * time.Second
	writeTimeout      = 600 * time.Second
)

func New(r http.Handler, port string) error {
	log.Info().Msgf("Starting server at port %s", port)

	srv := &http.Server{
		Handler:           r,
		Addr:              port,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		IdleTimeout:       idleTimeout,
		WriteTimeout:      writeTimeout,
	}

	return srv.ListenAndServe()
}

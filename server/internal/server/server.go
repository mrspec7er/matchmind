package server

import (
	"net/http"
	"os"
	"time"
)

type Config struct {
}

func NewInstance(s *Config) *http.Server {

	server := &http.Server{
		Addr:         os.Getenv("PORT"),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

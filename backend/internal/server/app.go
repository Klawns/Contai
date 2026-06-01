package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	config config
}

func NewServer() (*Server, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	deps, err := newDependencies(cfg)
	if err != nil {
		return nil, err
	}

	router := gin.Default()
	router.Use(securityHeaders(isProduction()))
	registerRoutes(router, deps)

	return &Server{
		router: router,
		config: cfg,
	}, nil
}

func (server *Server) Run() error {
	return server.httpServer().ListenAndServe()
}

func (server *Server) httpServer() *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%s", server.config.port),
		Handler:           server.router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

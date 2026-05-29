package server

import (
	"fmt"

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

	deps := newDependencies()
	router := gin.Default()
	registerRoutes(router, deps)

	return &Server{
		router: router,
		config: cfg,
	}, nil
}

func (server *Server) Run() error {
	return server.router.Run(fmt.Sprintf(":%s", server.config.port))
}

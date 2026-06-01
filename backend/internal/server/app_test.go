package server

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestHTTPServerTimeouts(t *testing.T) {
	server := Server{
		router: gin.New(),
		config: config{port: "8081"},
	}

	httpServer := server.httpServer()

	if httpServer.ReadTimeout != 10*time.Second {
		t.Fatalf("expected read timeout 10s, got %s", httpServer.ReadTimeout)
	}
	if httpServer.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf("expected read header timeout 5s, got %s", httpServer.ReadHeaderTimeout)
	}
	if httpServer.WriteTimeout != 15*time.Second {
		t.Fatalf("expected write timeout 15s, got %s", httpServer.WriteTimeout)
	}
	if httpServer.IdleTimeout != 60*time.Second {
		t.Fatalf("expected idle timeout 60s, got %s", httpServer.IdleTimeout)
	}
}

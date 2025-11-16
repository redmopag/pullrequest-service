package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: handler,
		},
	}
}

func (server *Server) Start() {
	log.Println("starting server on " + server.httpServer.Addr)
	if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v\n", err)
	}
}

func (server *Server) GracefulShutdown() {
	log.Println("shutting down server...")
	shutdwonCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.httpServer.Shutdown(shutdwonCtx); err != nil {
		log.Printf("Server Shutdown Failed:%+v", err)
	}
	log.Println("server Exited Properly")
}

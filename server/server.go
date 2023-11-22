package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	port int          // Port to listen on
	srv  *http.Server // HTTP server
}

func CreateServer(port int) *Server {
	// Create server
	server := &Server{
		port: port,
	}

	// Create HTTP server
	server.srv = &http.Server{
		Addr: fmt.Sprintf(":%d", server.port),
	}

	// Add handlers
	http.HandleFunc("/api/analyze_spec", server.analyzeSpecHandler())
	http.HandleFunc("/api/get_semantic_tokens", server.getSemanticTokensHandler())

	return server
}

func (server *Server) Serve() error {
	fmt.Printf("Starting server on :%d\n", server.port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		if err := server.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server stopped unexpectedly: %s\n", err.Error())
			stop <- os.Interrupt
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	<-stop
	fmt.Println("\nShutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown failed: %s\n", err.Error())
		return err
	}

	fmt.Println("Server gracefully stopped!")
	return nil
}

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ConfigMate/configmate/analyzer"
)

type Server struct {
	port     int               // Port to listen on
	analyzer analyzer.Analyzer // Analyzer to use for checking

	srv *http.Server // HTTP server
}

func CreateServer(port int, analyzer analyzer.Analyzer) *Server {
	// Create server
	server := &Server{
		port:     port,
		analyzer: analyzer,
	}

	// Create HTTP server
	server.srv = &http.Server{
		Addr: fmt.Sprintf(":%d", server.port),
	}

	// Add handlers
	http.HandleFunc("/api/check", server.checkHandler())

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

// checkHandler returns a handler for the check endpoint.
func (server *Server) checkHandler() http.HandlerFunc {
	// Return handler for check endpoint
	return func(w http.ResponseWriter, r *http.Request) {
		var rb analyzer.Rulebook

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&rb); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := server.analyzer.Analyze(rb)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

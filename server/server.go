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
	"github.com/ConfigMate/configmate/parsers"
	"github.com/ConfigMate/configmate/utils"
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
		var p struct {
			RulebookPath string `json:"rulebook_path"` // Path to the rulebook
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Read the rulebook file
		ruleBookData, err := os.ReadFile(p.RulebookPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Decode TOML into a Rulebook object
		ruleBook, err := utils.DecodeRulebook(ruleBookData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Parse rulebooks
		files := make(map[string]*parsers.Node)
		for alias, file := range ruleBook.Files {
			// Read the file
			data, err := os.ReadFile(file.Path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Parse the file
			parsedFile, err := parsers.Parse(data, file.Format)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Append the parse result to the files map
			files[alias] = parsedFile
		}

		// Get analyzer
		analyzer := &analyzer.AnalyzerImpl{}
		res, err := analyzer.AnalyzeConfigFiles(files, ruleBook.Rules)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

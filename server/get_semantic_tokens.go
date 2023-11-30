package server

import (
	"encoding/json"
	"net/http"

	"github.com/ConfigMate/configmate/files"
	"github.com/ConfigMate/configmate/langsrv"
)

type GetSemanticTokensRequest struct {
	Path    string `json:"path"`
	Content []byte `json:"content"`
}

type GetSemanticTokensResponse struct {
	SemanticTokens []langsrv.ParsedToken `json:"semantic_tokens"`
	Error          string                `json:"error"`
}

// getSemanticTokensHandler returns a handler for the get_semantic_tokens endpoint.
func (server *Server) getSemanticTokensHandler() http.HandlerFunc {
	// Return handler for check endpoint
	return func(w http.ResponseWriter, r *http.Request) {
		var p GetSemanticTokensRequest

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if p.Path != "" {
			// Read file content
			content, err := files.NewFileFetcher().FetchFile(p.Path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			p.Content = content
		}

		// Semantic tokens provider
		stp := langsrv.NewSemanticTokenProvider()
		tokens, err := stp.GetSemanticTokens(p.Content)
		errMessage := ""
		if err != nil {
			errMessage = err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&GetSemanticTokensResponse{
			SemanticTokens: tokens,
			Error:          errMessage,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

package server

import (
	"encoding/json"
	"net/http"

	"github.com/ConfigMate/configmate/analyzer"
	"github.com/ConfigMate/configmate/analyzer/check"
	"github.com/ConfigMate/configmate/analyzer/spec"
	"github.com/ConfigMate/configmate/files"
	"github.com/ConfigMate/configmate/parsers"
)

type AnalyzeSpecRequest struct {
	SpecFilePath string `json:"spec_file_path"`
}

type AnalyzeSpecResponse struct {
	Spec         *spec.Specification    `json:"spec"`
	CheckResults []analyzer.CheckResult `json:"check_results"`
	SpecError    *analyzer.SpecError    `json:"spec_error"`
}

// checkHandler returns a handler for the check endpoint.
func (server *Server) analyzeSpecHandler() http.HandlerFunc {
	// Return handler for check endpoint
	return func(w http.ResponseWriter, r *http.Request) {
		var p AnalyzeSpecRequest

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get analyzer
		a := analyzer.NewAnalyzer(
			spec.NewSpecParser(),
			check.NewCheckEvaluator(),
			files.NewFileFetcher(),
			parsers.NewParserProvider(),
		)

		spec, res, specError := a.AnalyzeSpecification(p.SpecFilePath)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&AnalyzeSpecResponse{
			Spec:         spec,
			CheckResults: res,
			SpecError:    specError,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

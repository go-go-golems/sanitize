package server

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strconv"
	"time"

	"github.com/go-go-golems/sanitize/examples"
	jsonsanitize "github.com/go-go-golems/sanitize/pkg/json"
	yamlsanitize "github.com/go-go-golems/sanitize/pkg/yaml"
)

//go:embed static
var staticFiles embed.FS

const (
	DefaultPort         = 8080
	maxRequestBodyBytes = 1 << 20
	readHeaderTimeout   = 5 * time.Second
	readTimeout         = 10 * time.Second
	writeTimeout        = 30 * time.Second
	idleTimeout         = 60 * time.Second
)

type analyzeRequest struct {
	Format string `json:"format"`
	Input  string `json:"input"`
}

type exampleResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Input       string `json:"input"`
	Category    string `json:"category,omitempty"`
	Source      string `json:"source,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Format      string `json:"format"`
}

func Run(port int) error {
	srv, err := New(port)
	if err != nil {
		return err
	}

	log.Info().Msg("listening")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func New(port int) (*http.Server, error) {
	if err := validatePort(port); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	// Serve embedded static files.
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, err
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// API: list examples.
	mux.HandleFunc("/api/examples", examplesHandler)

	// API: parse + sanitize.
	mux.HandleFunc("/api/sanitize", sanitizeHandler)

	// API: parse only (tree + errors, no fixing).
	mux.HandleFunc("/api/parse", parseHandler)

	return &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}, nil
}

func examplesHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONHeaders(w)

	var all []exampleResponse
	for _, ex := range yamlsanitize.Examples {
		all = append(all, exampleResponse{
			Name:        ex.Name,
			Description: ex.Description,
			Input:       ex.YAML,
			Category:    "demo",
			Source:      "builtin",
			Format:      "yaml",
		})
	}
	for _, ex := range examples.LoadFileExamples() {
		all = append(all, exampleResponse{
			Name:        ex.Name,
			Description: ex.Description,
			Input:       ex.YAML,
			Category:    ex.Category,
			Source:      ex.Source,
			Filename:    ex.Filename,
			Format:      "yaml",
		})
	}
	for _, ex := range jsonsanitize.Examples {
		all = append(all, exampleResponse{
			Name:        ex.Name,
			Description: ex.Description,
			Input:       ex.JSON,
			Category:    "demo",
			Source:      "builtin",
			Format:      "json",
		})
	}
	for _, ex := range examples.LoadJSONExamples() {
		all = append(all, exampleResponse{
			Name:        ex.Name,
			Description: ex.Description,
			Input:       ex.JSON,
			Category:    ex.Category,
			Source:      ex.Source,
			Filename:    ex.Filename,
			Format:      "json",
		})
	}

	if err := json.NewEncoder(w).Encode(all); err != nil {
		http.Error(w, "failed to encode examples", http.StatusInternalServerError)
	}
}

func sanitizeHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONHeaders(w)

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req analyzeRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	format, err := normalizeFormat(req.Format)
	if err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	switch format {
	case "yaml":
		if err := json.NewEncoder(w).Encode(yamlsanitize.Sanitize(req.Input)); err != nil {
			http.Error(w, "failed to encode sanitize result", http.StatusInternalServerError)
			return
		}
	case "json":
		if err := json.NewEncoder(w).Encode(jsonsanitize.Sanitize(req.Input)); err != nil {
			http.Error(w, "failed to encode sanitize result", http.StatusInternalServerError)
			return
		}
	}
}

func parseHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONHeaders(w)

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req analyzeRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	format, err := normalizeFormat(req.Format)
	if err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	switch format {
	case "yaml":
		treeText, errors, err := yamlsanitize.ParseTree(req.Input)
		if err != nil {
			http.Error(w, "parse error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"format":    "yaml",
			"tree_text": treeText,
			"errors":    errors,
		}); err != nil {
			http.Error(w, "failed to encode parse result", http.StatusInternalServerError)
		}
	case "json":
		treeText, errors, err := jsonsanitize.ParseTree(req.Input)
		if err != nil {
			http.Error(w, "parse error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"format":             "json",
			"tree_text":          treeText,
			"errors":             errors,
			"strict_parse_clean": jsonsanitize.StrictParse(req.Input) == nil,
		}); err != nil {
			http.Error(w, "failed to encode parse result", http.StatusInternalServerError)
		}
	}
}

func writeJSONHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}

	var extra json.RawMessage
	if err := dec.Decode(&extra); errors.Is(err, io.EOF) {
		return nil
	} else if err != nil {
		return err
	}
	if len(extra) > 0 {
		return errors.New("request body must contain a single JSON object")
	}
	return nil
}

func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid port %d: out of range", port)
	}
	return nil
}

func normalizeFormat(format string) (string, error) {
	switch format {
	case "", "yaml":
		return "yaml", nil
	case "json":
		return "json", nil
	default:
		return "", fmt.Errorf("unsupported format %q", format)
	}
}

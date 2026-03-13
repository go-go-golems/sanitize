package server

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-go-golems/sanitize/examples"
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

type yamlRequest struct {
	YAML string `json:"yaml"`
}

func Run(port int) error {
	srv, err := New(port)
	if err != nil {
		return err
	}

	log.Print("Listening")
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

	// Merge built-in examples (concise demos) with file-based corpus.
	var all []yamlsanitize.Example
	for _, ex := range yamlsanitize.Examples {
		ex.Source = "builtin"
		ex.Category = "demo"
		all = append(all, ex)
	}
	all = append(all, examples.LoadFileExamples()...)

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

	var req yamlRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	result := yamlsanitize.Sanitize(req.YAML)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode sanitize result", http.StatusInternalServerError)
		return
	}
}

func parseHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONHeaders(w)

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req yamlRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	treeText, errors, err := yamlsanitize.ParseTree(req.YAML)
	if err != nil {
		http.Error(w, "parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"tree_text": treeText,
		"errors":    errors,
	}); err != nil {
		http.Error(w, "failed to encode parse result", http.StatusInternalServerError)
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
	if dec.More() {
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

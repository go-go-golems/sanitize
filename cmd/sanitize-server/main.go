package main

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	yamlsanitize "github.com/go-go-golems/sanitize/pkg/yaml"
)

//go:embed static
var staticFiles embed.FS

const (
	maxRequestBodyBytes = 1 << 20
	readHeaderTimeout   = 5 * time.Second
	readTimeout         = 10 * time.Second
	writeTimeout        = 30 * time.Second
	idleTimeout         = 60 * time.Second
)

type yamlRequest struct {
	YAML string `json:"yaml"`
}

func main() {
	port, err := serverPort(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	srv, err := newServer(port)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Listening")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func newServer(port int) (*http.Server, error) {
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
	if err := json.NewEncoder(w).Encode(yamlsanitize.Examples); err != nil {
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

func serverPort(raw string) (int, error) {
	if raw == "" {
		raw = "8080"
	}

	port, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid PORT %q: %w", raw, err)
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("invalid PORT %q: out of range", raw)
	}
	return port, nil
}

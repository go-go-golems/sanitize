package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"

	yamlsanitize "github.com/go-go-golems/sanitize/pkg/yaml"
)

//go:embed static
var staticFiles embed.FS

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Serve embedded static files.
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// API: list examples.
	mux.HandleFunc("/api/examples", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		_ = json.NewEncoder(w).Encode(yamlsanitize.Examples)
	})

	// API: parse + sanitize.
	mux.HandleFunc("/api/sanitize", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

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

		var req struct {
			YAML string `json:"yaml"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
			return
		}

		result := yamlsanitize.Sanitize(req.YAML)
		_ = json.NewEncoder(w).Encode(result)
	})

	// API: parse only (tree + errors, no fixing).
	mux.HandleFunc("/api/parse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			YAML string `json:"yaml"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
			return
		}

		treeText, errors, err := yamlsanitize.ParseTree(req.YAML)
		if err != nil {
			http.Error(w, "parse error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"tree_text": treeText,
			"errors":    errors,
		})
	})

	log.Printf("Listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

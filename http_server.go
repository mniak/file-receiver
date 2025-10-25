package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mniak/file-receiver/internal/templates"
)

type HTTPServer struct {
	Port int

	server *http.Server
}

func (h *HTTPServer) Start() error {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/upload", h.rootEndpoint)
	h.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", h.Port),
		Handler: httpMux,
	}
	go func() {
		log.Printf("Starting server on port %d", h.Port)
		h.server.ListenAndServe()
	}()
	return nil
}

func (h *HTTPServer) Stop(ctx context.Context) error {
	if h.server != nil {
		return h.server.Shutdown(ctx)
	}
	return nil
}

const UploadDir = "UploadedFiles"

func (h *HTTPServer) rootEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := templates.Upload.Execute(w, nil); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodPost:

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao ler arquivo: %v", err), http.StatusBadRequest)
			return
		}
		filename := filepath.Base(header.Filename)
		filename = strings.TrimPrefix(filename, "/")
		filename = strings.TrimPrefix(filename, "\\")
		filename = strings.TrimPrefix(filename, ".")
		filename = strings.TrimSpace(filename)
		if len(filename) == 0 {
			http.Error(w, fmt.Sprintf("Nome do arquivo invalido: %q", filename), http.StatusInternalServerError)
			return
		}
		filename = filepath.Join(UploadDir, header.Filename)

		out, err := os.Create(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao abrir arquivo: %v", err), http.StatusInternalServerError)
			return
		}
		defer out.Close()
		if _, err = io.Copy(out, file); err != nil {
			http.Error(w, fmt.Sprintf("Erro ao escrever arquivo: %v", err), http.StatusInternalServerError)
			return
		}
		if err := templates.Uploaded.Execute(w, nil); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

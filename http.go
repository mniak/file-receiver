package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"

	"github.com/mniak/file-receiver/internal/templates"
)

type HTTPServer struct {
	Port int

	server *http.Server
}

func (h *HTTPServer) Start() error {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", h.rootEndpoint)
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

func (h *HTTPServer) rootEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if err := templates.Upload.Execute(w, nil); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

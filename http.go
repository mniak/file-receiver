package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
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
	var err error
	go func() {
		log.Printf("Starting server on port %d", h.Port)
		err = h.server.ListenAndServe()
	}()
	return err
}

func (h *HTTPServer) Stop() error {
	stopCtx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	if h.server != nil {
		return h.server.Shutdown(stopCtx)
	}
	return nil
}

func (h *HTTPServer) rootEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello File Receiver"))
}

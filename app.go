package receivefiles

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type App struct {
	Port             int
	ReceivedFilesDir string

	server *http.Server
	finish <-chan error
}

func (srv *App) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.uploadFormHandler)
	mux.HandleFunc("/submit", srv.fileUploadHandler)
	srv.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", srv.Port),
		Handler: mux,
	}
	finishChan := make(chan error)
	srv.finish = finishChan
	go func() {
		srv.server.ListenAndServe()
		close(finishChan)
	}()
}

func (srv *App) Wait() error {
	return <-srv.finish
}

func (srv *App) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.server.Shutdown(ctx)
}

package main

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"receivefiles/tmpl"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type HttpParams struct {
	Port             int
	ReceivedFilesDir string
}
type HttpService struct {
	HttpParams
	server *http.Server
	finish <-chan error
}

func (s *HttpService) Start() error {
	os.MkdirAll(s.ReceivedFilesDir, os.ModePerm)

	r := gin.Default()
	r.LoadHTMLFS(http.FS(tmpl.FS()), "*.html")

	r.Use(s.onError)
	r.NoRoute(s.onNotFound)
	staticFS, err := fs.Sub(StaticFS, "static")
	if err != nil {
		return errors.WithMessage(err, "failed to find static/ in the embedded fs")
	}
	r.StaticFS("/static", http.FS(staticFS))

	r.GET("/", s.getRoot)
	r.POST("/submit", s.postSubmit)

	addr := fmt.Sprintf(":%d", s.Port)
	finishChan := make(chan error)
	s.server = &http.Server{
		Addr:    addr,
		Handler: r,
	}
	s.finish = finishChan

	go func() {
		s.server.ListenAndServe()
		close(finishChan)
	}()
	return nil
}

func (s *HttpService) Wait() error {
	return <-s.finish
}

func (s *HttpService) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *HttpService) onNotFound(c *gin.Context) {
	c.JSON(404, "not found")
}

const qSuccessMessage = "sm"

func (s *HttpService) getRoot(c *gin.Context) {
	msg := c.Query(qSuccessMessage)
	c.HTML(200, "index.html", map[string]any{
		"SuccessText": msg,
	})
}

func (s *HttpService) postSubmit(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		c.Error(errors.WithMessage(err, "failed to get file"))
		return
	}
	dstPath := filepath.Join(s.ReceivedFilesDir, f.Filename)
	err = c.SaveUploadedFile(f, dstPath)
	if err != nil {
		c.Error(errors.WithMessage(err, "failed to save file"))
		return

	}
	q := make(url.Values)
	q.Add(qSuccessMessage, fmt.Sprintf("Arquivo '%s' enviado!", f.Filename))
	c.Redirect(302, "/?"+q.Encode())
}

func (s *HttpService) onError(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		err := c.Errors.Last().Err
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": err.Error(),
		})
	}
}

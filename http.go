package receivefiles

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
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
	os.MkdirAll(s.HttpParams.ReceivedFilesDir, os.ModePerm)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.LoadHTMLFS(http.FS(tmpl.FS()), "*.html")
	r.GET("/", s.uploadFormHandler)
	r.POST("/submit", s.fileUploadHandler)

	staticFS, err := fs.Sub(StaticFS, "static")
	if err != nil {
		return errors.WithMessage(err, "failed to find /static in the embedded filesystem")
	}
	r.StaticFS("/static", http.FS(staticFS))
	r.NoRoute(s.notFound)

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

func (s *HttpService) notFound(c *gin.Context) {
	c.JSON(404, "not found")
}

func (s *HttpService) uploadFormHandler(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func (s *HttpService) fileUploadHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.Error(errors.WithMessage(err, "parse form"))
		return
	}
	files, ok := form.File["file"]
	if !ok {
		c.Error(errors.New("missing file"))
		return
	}

	for _, f := range files {
		src, err := f.Open()
		if err != nil {
			c.Error(errors.WithMessage(err, "open file sent"))
			return
		}

		dstPath := filepath.Join(s.ReceivedFilesDir, f.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			c.Error(errors.WithMessage(err, "create file on server"))
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			c.Error(errors.WithMessage(err, "write file"))
			return
		}

	}
	c.JSON(200, "File uploaded")
}

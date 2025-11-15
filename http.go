package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Server struct {
	Flags Flags
}

func (srv *Server) run() {
	http.HandleFunc("/", srv.uploadFormHandler)
	http.HandleFunc("/submit", srv.fileUploadHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", srv.Flags.Port), nil)
}

func (srv *Server) uploadFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl := `
    <html>
    <body>
        <h1>IASD Cidade das Flores</h1>
		<h2>Enviar arquivo</h2>
        <form action="/submit" method="post" enctype="multipart/form-data">
            <input type="file" name="file" />
            <input type="submit" value="Enviar" />
        </form>
    </body>
    </html>`

	t := template.Must(template.New("upload").Parse(tmpl))
	t.Execute(w, nil)
}

// POST /upload/submit - Handles file upload
func (srv *Server) fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // Limit: 10MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get file from posted form-data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create destination file
	dstPath := filepath.Join(srv.Flags.ReceivedFilesDir, handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file contents
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", handler.Filename)
}

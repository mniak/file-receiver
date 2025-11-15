package receivefiles

import (
	"fmt"
	"net/http"
)

type App struct {
	Port             int
	ReceivedFilesDir string
}

func (srv *App) Run() {
	http.HandleFunc("/", srv.uploadFormHandler)
	http.HandleFunc("/submit", srv.fileUploadHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", srv.Port), nil)
}

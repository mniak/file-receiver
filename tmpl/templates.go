package tmpl

import (
	"embed"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed *.html
var _fs embed.FS

func FS() embed.FS {
	return _fs
}

func getTemplate(name string) (*template.Template, error) {
	t, err := template.ParseFS(_fs, fmt.Sprintf("%s.html", name))
	return t, err
}

func Index() *template.Template {
	return template.Must(getTemplate("index"))
}

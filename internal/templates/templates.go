package templates

import (
	_ "embed"
	"html/template"
)

//go:embed upload.html
var upload string
var Upload = template.Must(template.New("upload").Parse(upload))

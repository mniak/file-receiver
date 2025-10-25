package templates

import (
	_ "embed"
	"html/template"
)

//go:embed upload.html
var upload string
var Upload = template.Must(template.New("upload").Parse(upload))

//go:embed uploaded.html
var uploaded string
var Uploaded = template.Must(template.New("uploaded").Parse(uploaded))

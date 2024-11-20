package assets

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

//go:embed templates/*.html
var TemplatesFS embed.FS

//go:embed static/*
var StaticFS embed.FS

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Template {
	return &Template{templates: template.Must(template.ParseFS(TemplatesFS, "templates/*.html"))}
}

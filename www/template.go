package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

// Template app template
type Template struct {
	templates *template.Template
}

// Render echo.Renderer interface implementation
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

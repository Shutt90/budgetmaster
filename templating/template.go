package template

import (
	"io"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func NewTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func RenderNegative(message string, negative bool): string {
	return `<div class="justify-center w-screen py-2 absolute top-0 left-0">
    {{ template "flash" . }}
  </div>`
}

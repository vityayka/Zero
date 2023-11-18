package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/vityayka/go-zero/context"
	"github.com/vityayka/go-zero/models"
)

func Must(tpl Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return tpl
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(patterns[0])

	tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField has never been implemented :(")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser has never been implemented :(")
			},
		},
	)

	tpl, err := tpl.ParseFS(fs, patterns...)

	// var ErrorNotFound = errors.New("not found")
	// log.Println(errors.Is(ErrorNotFound, ErrorNotFound))

	if err != nil {
		log.Printf("error parsing template %v", err)
		return Template{}, fmt.Errorf("error parsing template %v", err)
	}
	return Template{HTMLTpl: tpl}, nil
}

type Template struct {
	HTMLTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) bool {
	// log.Println("hey june")
	tpl, err := t.HTMLTpl.Clone()
	if err != nil {
		log.Printf("cloning template failed: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
	tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
		},
	)

	var buf bytes.Buffer

	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("error executing template %v", err)
		http.Error(w, "error executing template", http.StatusInternalServerError)
		return true
	}

	buf.WriteTo(w)
	buf.Reset()
	// makes

	return false
}

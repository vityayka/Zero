package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/csrf"
	"github.com/vityayka/go-zero/context"
	"github.com/vityayka/go-zero/models"
)

type public interface {
	Public() string
}

func Must(tpl Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return tpl
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(path.Base(patterns[0]))

	tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField has never been implemented :(")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser has never been implemented :(")
			},
			"errors": func() []string {
				return nil
			},
		},
	)

	tpl, err := tpl.ParseFS(fs, patterns...)

	if err != nil {
		log.Printf("error parsing template %v", err)
		return Template{}, fmt.Errorf("error parsing template %v", err)
	}
	return Template{HTMLTpl: tpl}, nil
}

type Template struct {
	HTMLTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) bool {
	tpl, err := t.HTMLTpl.Clone()
	if err != nil {
		log.Printf("cloning template failed: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
	errMsgs := errMessages(errs...)
	tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				return errMsgs
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

	return false
}

func errMessages(errs ...error) []string {
	var messages []string
	for _, err := range errs {
		var publicErr public
		if errors.As(err, &publicErr) {
			messages = append(messages, publicErr.Public())
		} else {
			messages = append(messages, "Something went wrong")
		}
	}
	return messages
}

package views

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

func Must(tpl Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return tpl
}

func ParseFS(fs fs.FS, pattern string) (Template, error) {
	tpl, err := template.ParseFS(fs, pattern)

	var ErrorNotFound = errors.New("not found")
	log.Println(errors.Is(ErrorNotFound, ErrorNotFound))

	if err != nil {
		log.Printf("error parsing template %v", err)
		return Template{}, fmt.Errorf("error parsing template %v", err)
	}
	return Template{tpl}, nil
}

func Parse(path string) (Template, error) {
	tpl, err := template.ParseFiles(path)

	var ErrorNotFound = errors.New("not found")
	log.Println(errors.Is(ErrorNotFound, ErrorNotFound))

	if err != nil {
		log.Printf("error parsing template %v", err)
		return Template{}, fmt.Errorf("error parsing template %v", err)
	}
	return Template{tpl}, nil
}

type Template struct {
	HTMLTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) bool {
	log.Println("hey june")
	err := t.HTMLTpl.Execute(w, data)
	if err != nil {
		log.Printf("error executing template %v", err)
		http.Error(w, "error executing template", http.StatusInternalServerError)
		return true
	}
	return false
}

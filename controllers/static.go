package controllers

import (
	"net/http"

	"github.com/vityayka/go-zero/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func Photos(tpl views.Template) http.HandlerFunc {

	questions := []struct {
		Path  string
		Title string
	}{
		{
			Path:  "/2022/e99fqiwef.jpg",
			Title: "Random shit",
		},
		{
			Path:  "/2022/f9a9refi4j23.jpg",
			Title: "Another rand shit",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}

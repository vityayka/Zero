package controllers

import (
	"fmt"
	"net/http"
)

type Users struct {
	Templates struct {
		New Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	// email := r.PostForm.Get("email")
	// password := r.PostForm.Get("password")
	file, fileHeader, err := r.FormFile("photo")
	fileInfo := []any{fileHeader.Filename, fileHeader.Size}
	fmt.Printf("file: %+v\n", file)
	fmt.Printf("info: %+v\n", fileInfo)
	fmt.Printf("err: %+v\n", err)
}

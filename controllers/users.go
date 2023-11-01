package controllers

import (
	"fmt"
	"net/http"

	"github.com/vityayka/go-zero/models"
)

type Users struct {
	Templates struct {
		New    Template
		Signin Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) Signin(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	u.Templates.Signin.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user, err := u.UserService.Create(email, password)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "Created user: %v", user)
}

func (u Users) Auth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user, err := u.UserService.Auth(email, password)

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}

	fmt.Fprintf(w, "Authenticated user: %v", user)
}

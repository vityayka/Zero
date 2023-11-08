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
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	sessionTokenCookie, err := r.Cookie(CookieSession)

	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/users/signin", http.StatusFound)
	}

	fmt.Fprintf(w, "session: %v", sessionTokenCookie.Value)
}

func (u Users) Signin(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	u.Templates.Signin.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user, err := u.UserService.Create(email, password)

	if err != nil {
		panic(err)
	}

	u.createSession(w, user)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) Auth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user, err := u.UserService.Auth(email, password)

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		fmt.Printf("Error: %v", err)
		return
	}

	u.createSession(w, user)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) createSession(w http.ResponseWriter, user *models.User) {
	session, err := u.SessionService.Create(user.ID)

	if err != nil {
		panic(err)
	}

	http.SetCookie(w, newCookie(CookieSession, session.Token))
}

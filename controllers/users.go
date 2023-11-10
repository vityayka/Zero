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

	if err != nil {
		http.Redirect(w, r, "/users/signin", http.StatusFound)
		return
	}

	user, err := u.SessionService.User(sessionTokenCookie.Value)
	if err != nil {
		fmt.Fprintf(w, "something went wrong: %v", err)
		return

	}
	fmt.Fprintf(w, "user: %v", user)
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

func (u Users) SignOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(CookieSession)

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		fmt.Printf("Error: %v", err)
		return
	}

	err = u.SessionService.Delete(cookie.Value)
	if err != nil {
		http.Error(w, "Smthng went wrong", http.StatusInternalServerError)
		fmt.Errorf("deleting session: %v", err)
		return
	}

	deleteCookie(w, CookieSession)

	http.Redirect(w, r, "/users/signin", http.StatusFound)
}

func (u Users) createSession(w http.ResponseWriter, user *models.User) {
	session, err := u.SessionService.Create(user.ID)

	if err != nil {
		panic(err)
	}

	http.SetCookie(w, newCookie(CookieSession, session.Token))
}
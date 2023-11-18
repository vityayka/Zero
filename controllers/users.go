package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/vityayka/go-zero/context"
	"github.com/vityayka/go-zero/models"
)

type Users struct {
	Templates struct {
		New            Template
		ForgotPassword Template
		NewPassword    Template
		Signin         Template
	}
	UserService       *models.UserService
	SessionService    *models.SessionService
	ResetTokenService *models.ResetTokenService
	EmailService      *models.EmailService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

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

func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.PostForm.Get("email")

	token, err := u.ResetTokenService.Create(email)

	if err != nil {
		http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
		return
	}

	user, err := u.ResetTokenService.User(token.Token)

	if err != nil {
		http.Error(w, "No such user :(", http.StatusNotFound)
		return
	}

	url := url.Values{
		"token": {token.Token},
	}

	u.EmailService.ForgotPassword(user.Email, "http://localhost:3000/users/reset-password?"+url.Encode())

	fmt.Fprintf(w, "Go to your email inbox")
}

func (u Users) ConsumeResetToken(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	user, err := u.ResetTokenService.Consume(token)
	if err != nil {
		http.Error(w, "provided token is bad", http.StatusUnauthorized)
	}

	var data struct {
		User *models.User
	}

	data.User = user

	u.Templates.NewPassword.Execute(w, r, data)
}

func (u Users) ProcessNewPassword(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.FormValue("user_id"))
	password := r.FormValue("password")
	passwordRepeat := r.FormValue("password_repeat")

	if password != passwordRepeat {
		http.Error(w, "Passwords don't match", http.StatusBadRequest)
	}

	u.UserService.UpdatePassword(userID, password)
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
		log.Fatalf("deleting session: %v", err)
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

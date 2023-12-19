package controllers

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/vityayka/go-zero/context"
	"github.com/vityayka/go-zero/models"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"strings"
)

type OAuth struct {
	ProviderConfigs map[string]*oauth2.Config
	TokenService    models.OAuthService
}

func (oa OAuth) Connect(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	provider = strings.ToLower(provider)
	cfg, isOk := oa.ProviderConfigs[provider]
	if !isOk {
		http.Error(w, "Invalid oauth provider", http.StatusBadRequest)
		return
	}

	state := csrf.Token(r)
	cookie := newCookie("oauth_state", state)
	http.SetCookie(w, cookie)

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := cfg.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("redirect_uri", redirectURI(r, provider)),
	)

	http.Redirect(w, r, url, http.StatusFound)
}

func (oa OAuth) Callback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	provider = strings.ToLower(provider)
	cfg, isOk := oa.ProviderConfigs[provider]
	if !isOk {
		http.Error(w, "Invalid oauth provider", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, "This is a fraud", http.StatusBadRequest)
		return
	}

	state := r.FormValue("state")
	if cookie.Value != state {
		http.Error(w, "This is a fraud", http.StatusBadRequest)
		return
	}

	deleteCookie(w, "oauth_state")

	code := r.FormValue("code")

	if cookie.Value != state {
		http.Error(w, "Unknown oauth error", http.StatusInternalServerError)
		log.Printf("cfg.Exchange errror: %v", err)
		return
	}

	token, err := cfg.Exchange(
		r.Context(),
		code,
		oauth2.SetAuthURLParam("redirect_uri", redirectURI(r, provider)),
	)

	if err != nil {
		var retrieveError *oauth2.RetrieveError
		if errors.As(err, &retrieveError) {
			text := string(retrieveError.Body)
			fmt.Fprintf(w, "exchange error: %v", text)
			return
		}
	}

	user := context.User(r.Context())
	_, err = oa.TokenService.Create(user.ID, token.AccessToken, token.Expiry, provider)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "token: %v", token)
}

func redirectURI(r *http.Request, provider string) string {
	var scheme string
	if strings.Contains(r.Host, "localhost") {
		scheme = "http"
	} else {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s/oauth/%s/callback", scheme, r.Host, provider)
}

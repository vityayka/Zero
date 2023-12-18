package controllers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

type OAuth struct {
	ProviderConfigs map[string]*oauth2.Config
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
	r.URL.Query()
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
		oauth2.SetAuthURLParam("redirect_uri", "http://localhost:3000/oauth/dropbox/callback"),
	)

	http.Redirect(w, r, url, http.StatusFound)
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

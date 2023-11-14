package middlewares

import (
	"fmt"
	"net/http"

	"github.com/vityayka/go-zero/context"
	"github.com/vityayka/go-zero/controllers"
	"github.com/vityayka/go-zero/models"
)

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (u UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionTokenCookie, err := r.Cookie(controllers.CookieSession)
		var user *models.User

		if err != nil {
			fmt.Printf("cookie err: %v", err)
			user = nil
		} else {
			user, err = u.SessionService.User(sessionTokenCookie.Value)
			if err != nil {
				fmt.Printf("something went wrong: %v", err)
				user = nil
			}
		}

		ctx := context.WithUser(r.Context(), user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (u UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/users/signin", http.StatusFound)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

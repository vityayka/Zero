package middlewares

import (
	"net/http"

	"github.com/gorilla/csrf"
)

type CSRFConfig struct {
	Key    string
	Secure bool
}

func CSRFProtect(cfg CSRFConfig) func(http.Handler) http.Handler {
	return csrf.Protect(
		[]byte(cfg.Key),
		csrf.Secure(cfg.Secure),
		csrf.Path("/"),
	)
}

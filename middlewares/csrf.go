package middlewares

import (
	"net/http"

	"github.com/gorilla/csrf"
)

const csrfKey string = "sf9gia0sdcvif04FF349fadvununaEEE"

func Protect() func(http.Handler) http.Handler {
	return csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false), // TODO: move it to .env or an analog
		csrf.Path("/"),
	)
}

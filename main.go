package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vityayka/go-zero/controllers"
	"github.com/vityayka/go-zero/models"
	"github.com/vityayka/go-zero/templates"
	"github.com/vityayka/go-zero/views"
)

func main() {
	router := chi.NewRouter()
	// router.Use(middleware.Logger)

	router.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "hello.gohtml", "tailwind.gohtml"))))
	router.Get("/dashboard",
		controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "dashboard.gohtml", "tailwind.gohtml"))))

	db, err := models.Open(models.DefaultPgConfig())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	UserService := models.UserService{
		DB: db,
	}
	SessionService := models.SessionService{
		DB: db,
	}
	usersC := controllers.Users{
		UserService:    &UserService,
		SessionService: &SessionService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.Signin = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))

	router.Get("/users/signup", usersC.New)
	// router.Get("/users/signin", TimeMiddleware(usersC.Signin))
	router.Get("/users/signin", usersC.Signin)
	router.Post("/users/new", usersC.Create)
	router.Post("/users/auth", usersC.Auth)

	router.Route("/photos", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/{photoSlug:[a-zA-z-0-9]+}",
			controllers.Photos(views.Must(views.ParseFS(templates.FS, "photos.gohtml", "tailwind.gohtml"))))
	})
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "PAGE NOT FOUND", http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	csrfKey := "sf9gia0sdcvif04FF349fadvununaEEE"
	csrfMiddleware := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false), // TODO: move it to .env or an analog
		csrf.Path("/"),
	)
	http.ListenAndServe(":3000", csrfMiddleware(router))
}

// func TimeMiddleware(h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		h(w, r)
// 		fmt.Fprintf(w, "Request time = %s", time.Since(start).String())
// 	}
// }

package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vityayka/go-zero/controllers"
	"github.com/vityayka/go-zero/middlewares"
	"github.com/vityayka/go-zero/migrations"
	"github.com/vityayka/go-zero/models"
	"github.com/vityayka/go-zero/templates"
	"github.com/vityayka/go-zero/views"
)

func main() {
	db := setupDB()
	defer db.Close()
	usersC := initControllers(db)

	router := setupRoutes(usersC, db)

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", router)
}

func setupRoutes(usersC controllers.Users, db *sql.DB) *chi.Mux {
	router := chi.NewRouter()

	userMiddleware := middlewares.UserMiddleware{
		SessionService: &models.SessionService{DB: db},
	}

	router.Use(middlewares.Protect(), userMiddleware.SetUser)

	router.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "hello.gohtml", "tailwind.gohtml"))))
	router.Get("/dashboard",
		controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "dashboard.gohtml", "tailwind.gohtml"))))

	router.Get("/users/signup", usersC.New)
	router.Get("/users/signin", usersC.Signin)
	router.Post("/users/new", usersC.Create)
	router.Post("/users/auth", usersC.Auth)
	router.Post("/users/signout", usersC.SignOut)
	router.Route("/users/me", func(r chi.Router) {
		r.Use(userMiddleware.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	router.Route("/photos", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/{photoSlug:[a-zA-z-0-9]+}",
			controllers.Photos(views.Must(views.ParseFS(templates.FS, "photos.gohtml", "tailwind.gohtml"))))
	})
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "PAGE NOT FOUND", http.StatusNotFound)
	})
	return router
}

func initControllers(db *sql.DB) controllers.Users {
	usersC := controllers.Users{
		UserService:    &models.UserService{DB: db},
		SessionService: &models.SessionService{DB: db},
	}

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.Signin = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	return usersC
}

func setupDB() *sql.DB {
	db, err := models.Open(models.DefaultPgConfig())
	if err != nil {
		panic(err)
	}

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	return db
}

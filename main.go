package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/lpernett/godotenv"
	"github.com/vityayka/go-zero/controllers"
	"github.com/vityayka/go-zero/middlewares"
	"github.com/vityayka/go-zero/migrations"
	"github.com/vityayka/go-zero/models"
	"github.com/vityayka/go-zero/templates"
	"github.com/vityayka/go-zero/views"
)

type config struct {
	SMTP   models.SMTPConfig
	CSRF   middlewares.CSRFConfig
	Server struct {
		Address string
	}
	Postgres models.PostgresConfig
}

func loadEnvConfig() (config, error) {
	var cfg config

	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env load failed: %v", err)
		panic(err)
	}

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Port, _ = strconv.Atoi(os.Getenv("SMTP_Port"))
	cfg.SMTP.User = os.Getenv("SMTP_USER")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	cfg.Postgres.Database = os.Getenv("POSTGRES_DATABASE")
	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.User = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.Sslmode = os.Getenv("POSTGRES_SSLMODE")

	return cfg, err
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	db := setupDB(cfg.Postgres)
	defer db.Close()
	usersC := initControllers(db, cfg)

	router := setupRoutes(usersC, db, cfg)

	fmt.Printf("Starting the server on %s...", cfg.Server.Address)
	http.ListenAndServe(cfg.Server.Address, router)
}

func setupRoutes(usersC controllers.Users, db *sql.DB, cfg config) *chi.Mux {
	router := chi.NewRouter()

	userMiddleware := middlewares.UserMiddleware{
		SessionService: &models.SessionService{DB: db},
	}

	router.Use(middlewares.CSRFProtect(cfg.CSRF), userMiddleware.SetUser)

	router.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "hello.gohtml", "tailwind.gohtml"))))
	router.Get("/dashboard",
		controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "dashboard.gohtml", "tailwind.gohtml"))))

	router.Get("/users/signup", usersC.New)
	router.Get("/users/signin", usersC.Signin)
	router.Post("/users/new", usersC.Create)
	router.Post("/users/auth", usersC.Auth)
	router.Post("/users/signout", usersC.SignOut)
	router.Get("/users/forgot-password", usersC.ForgotPassword)
	router.Post("/users/forgot-password", usersC.ProcessForgotPassword)
	router.Get("/users/reset-password", usersC.NewPassword)
	router.Post("/users/new-password", usersC.ProcessNewPassword)
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

func initControllers(db *sql.DB, cfg config) controllers.Users {
	userService := &models.UserService{DB: db}
	usersC := controllers.Users{
		UserService:       userService,
		SessionService:    &models.SessionService{DB: db},
		ResetTokenService: &models.PasswordResetService{DB: db, UserService: userService, Duration: time.Hour},
		EmailService:      models.NewEmailService(cfg.SMTP),
	}

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.Signin = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "forgot_password.gohtml", "tailwind.gohtml"))
	usersC.Templates.NewPassword = views.Must(views.ParseFS(templates.FS, "new_password.gohtml", "tailwind.gohtml"))
	return usersC
}

func setupDB(cfg models.PostgresConfig) *sql.DB {
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	return db
}

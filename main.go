package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vityayka/go-zero/controllers"
	"github.com/vityayka/go-zero/templates"
	"github.com/vityayka/go-zero/views"
)

func main() {
	// azaz := []struct {
	// 	X string
	// 	Y int
	// }{{"azaz", 1}, {"fdfdf", 2}}

	// fmt.Printf("Azaz: %+v", azaz)

	router := chi.NewRouter()
	// router.Use(middleware.Logger)

	router.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "hello.gohtml"))))
	router.Get("/dashboard",
		controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "dashboard.gohtml"))))

	router.Route("/photos", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/{photoSlug:[a-zA-z-0-9]+}",
			controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "custom.gohtml"))))
	})
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "PAGE NOT FOUND", http.StatusNotFound)
	})
	time.Sleep(1 * time.Second)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", router)
}

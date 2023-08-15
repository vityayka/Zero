package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vityayka/go-zero/views"
)

func execTpl(w http.ResponseWriter, path string, data any) bool {
	tpl, err := views.Parse(path)
	if err != nil {
		http.Error(w, "error parsing template", http.StatusInternalServerError)
	}
	return tpl.Execute(w, nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	execTpl(w, filepath.Join("templates", "hello.gohtml"), nil)
}

func handlerDash(w http.ResponseWriter, r *http.Request) {
	execTpl(w, filepath.Join("templates", "dashboard.gohtml"), nil)
}

func handleCustom(w http.ResponseWriter, r *http.Request) {
	execTpl(w, filepath.Join("templates", "custom.gohtml"), nil)
}

func main() {
	// azaz := []struct {
	// 	X string
	// 	Y int
	// }{{"azaz", 1}, {"fdfdf", 2}}

	// fmt.Printf("Azaz: %+v", azaz)

	router := chi.NewRouter()
	// router.Use(middleware.Logger)
	router.Get("/", handleRoot)
	router.Get("/dashboard", handlerDash)
	router.Route("/photos", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/{photoSlug:[a-zA-z-0-9]+}", handleCustom)
	})
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "PAGE NOT FOUND", http.StatusNotFound)
	})
	time.Sleep(1 * time.Second)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", router)
}

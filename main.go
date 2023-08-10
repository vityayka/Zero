package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func execTpl(w http.ResponseWriter, tpl string) bool {
	w.Header().Set("content-type", "text/html")
	t, err := template.ParseFiles(tpl)
	if err != nil {
		log.Printf("error parsing template %v", err)
		http.Error(w, "error parsing template", http.StatusInternalServerError)
		return true
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("error executing template %v", err)
		http.Error(w, "error executing template", http.StatusInternalServerError)
		return true
	}
	return false
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	execTpl(w, filepath.Join("templates", "hello.gohtml"))
}

func handlerDash(w http.ResponseWriter, r *http.Request) {
	execTpl(w, filepath.Join("templates", "dashboard.gohtml"))
}

func handleCustom(w http.ResponseWriter, r *http.Request) {
	execTpl(w, filepath.Join("templates", "custom.gohtml"))
}

func main() {
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

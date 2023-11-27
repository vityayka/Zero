package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/vityayka/go-zero/context"
	"github.com/vityayka/go-zero/models"
)

type Galleries struct {
	Templates struct {
		New  Template
		Edit Template
	}
	Service *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Templates.New.Execute(w, r, data)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")

	gallery, err := g.Service.Create(title, context.User(r.Context()).ID)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)

	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}
	gallery, err := g.Service.ById(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	user := context.User(r.Context())
	if user.ID != gallery.UserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	g.Templates.Edit.Execute(w, r, gallery)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}

	gallery, err := g.Service.ById(id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	gallery.Title = r.FormValue("title")

	err = g.Service.Update(gallery)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)

	http.Redirect(w, r, editPath, http.StatusFound)
}

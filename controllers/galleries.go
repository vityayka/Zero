package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/vityayka/go-zero/context"
	"github.com/vityayka/go-zero/models"
)

type Galleries struct {
	Templates struct {
		Show  Template
		New   Template
		Edit  Template
		Index Template
	}
	Service *models.GalleryService
}

type GalleryOutput struct {
	ID    int
	Title string
}

func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Galleries []GalleryOutput
	}

	user := context.User(r.Context())

	galleries, err := g.Service.ByUserId(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, GalleryOutput{
			ID:    gallery.ID,
			Title: gallery.Title,
		})
	}

	g.Templates.Index.Execute(w, r, data)
}

func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Gallery GalleryOutput
		Images  []struct {
			URL string
		}
	}

	gallery, err := g.galleryById(r, w)
	if err != nil {
		return
	}

	data.Gallery = GalleryOutput{
		ID:    gallery.ID,
		Title: gallery.Title,
	}

	for i := 0; i < 10; i++ {
		w, h := rand.Intn(500)+200, rand.Intn(500)+200
		imgURL := fmt.Sprintf("https://placekitten.com/%d/%d", w, h)
		var image struct {
			URL string
		}
		image.URL = imgURL
		data.Images = append(data.Images, image)
	}

	g.Templates.Show.Execute(w, r, data)
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
	gallery, err := g.galleryById(r, w, galleryBelongsToUser)
	if err != nil {
		return
	}

	g.Templates.Edit.Execute(w, r, gallery)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(r, w, galleryBelongsToUser)
	if err != nil {
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

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(r, w, galleryBelongsToUser)
	if err != nil {
		return
	}

	err = g.Service.Delete(gallery)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/galleries", http.StatusFound)
}

type galleryCriteria func(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error

func (g Galleries) galleryById(r *http.Request, w http.ResponseWriter, opts ...galleryCriteria) (*models.Gallery, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Provided id is invalid", http.StatusBadRequest)
		return nil, err
	}

	gallery, err := g.Service.ById(id)
	if err != nil {
		http.Error(w, "Gallery is not found", http.StatusNotFound)
		return nil, err
	}

	for _, opt := range opts {
		err = opt(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}

	return gallery, nil
}

func galleryBelongsToUser(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())
	if user.ID != gallery.UserID {
		http.Error(w, "Gallery is not found", http.StatusNotFound)
		return models.ErrUnauthorized
	}
	return nil
}

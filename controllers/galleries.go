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
		http.Error(w, err.Error(), models.HttpErrorCode(err))
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
	gallery, err := g.galleryById(r, w)
	if err != nil {
		http.Error(w, err.Error(), models.HttpErrorCode(err))
		return
	}

	g.Templates.Edit.Execute(w, r, gallery)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(r, w)
	if err != nil {
		http.Error(w, err.Error(), models.HttpErrorCode(err))
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
	gallery, err := g.galleryById(r, w)
	if err != nil {
		http.Error(w, err.Error(), models.HttpErrorCode(err))
		return
	}

	err = g.Service.Delete(gallery)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g Galleries) galleryById(r *http.Request, w http.ResponseWriter) (*models.Gallery, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return nil, models.ErrBadRequest
	}

	gallery, err := g.Service.ById(id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	user := context.User(r.Context())
	if user.ID != gallery.UserID {
		return nil, models.ErrUnauthorized
	}
	return gallery, nil
}

package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vityayka/go-zero/context"
	"github.com/vityayka/go-zero/models"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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

type ImageOutput struct {
	GalleryID       int
	Filename        string
	FilenameEscaped string
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
		Images  []ImageOutput
	}

	gallery, err := g.galleryById(r, w)
	if err != nil {
		return
	}

	data.Gallery = GalleryOutput{
		ID:    gallery.ID,
		Title: gallery.Title,
	}

	outputImages, err := g.outputImages(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	data.Images = outputImages

	g.Templates.Show.Execute(w, r, data)
}

func (g Galleries) outputImages(gallery *models.Gallery) ([]ImageOutput, error) {
	images, err := g.Service.Images(gallery.ID)
	if err != nil {
		return nil, err
	}

	var outputImages []ImageOutput
	for _, image := range images {
		outputImages = append(outputImages, ImageOutput{
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
			GalleryID:       image.GalleryID,
		})
	}
	return outputImages, nil
}

func (g Galleries) Image(w http.ResponseWriter, r *http.Request) {
	galleryId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid gallery id", http.StatusBadRequest)
		return
	}

	filename := chi.URLParam(r, "filename")

	image, err := g.Service.Image(galleryId, filename)
	if err != nil {
		http.Error(w, "Image is not found", http.StatusNotFound)
	}

	http.ServeFile(w, r, image.Path)
}

func (g Galleries) UploadImages(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(r, w, galleryBelongsToUser)
	if err != nil {
		return
	}

	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		log.Printf("uploading an image to a gallery: %v", err)
		http.Error(w, "Uploading failed", http.StatusInternalServerError)
		return
	}

	fileHeaders := r.MultipartForm.File["images"]

	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Fail", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		err = g.Service.CreateImage(fileHeader.Filename, gallery.ID, file)
		if err != nil {
			var fileError models.FileError
			if errors.As(err, &fileError) {
				msg := fmt.Sprintf("%s has an invalid content type or extension. Only jpeg, gif and png images "+
					"are supported", fileHeader.Filename)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			http.Error(w, "Uploading failed", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/galleries/%d", gallery.ID), http.StatusFound)
}

func (g Galleries) UploadExternalImages(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(r, w, galleryBelongsToUser)
	if err != nil {
		http.Error(w, "Fail", http.StatusInternalServerError)
		return
	}

	var links struct {
		Links []string `json:"links"`
	}

	err = json.NewDecoder(r.Body).Decode(&links)
	if err != nil {
		http.Error(w, "Fail", http.StatusInternalServerError)
		return
	}

	var eg errgroup.Group

	for _, imageUrl := range links.Links {
		insideLoopUrl := imageUrl
		fmt.Println(insideLoopUrl)
		eg.Go(func() error {
			return g.Service.CreateImageFromUrl(gallery.ID, insideLoopUrl)
		})
	}
	err = eg.Wait()
	if err != nil {
		http.Error(w, "Unable to upload some of images", http.StatusInternalServerError)
		return
	}

	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(r, w, galleryBelongsToUser)

	filename := filepath.Base(chi.URLParam(r, "filename"))

	image, err := g.Service.Image(gallery.ID, filename)
	if err != nil {
		http.Error(w, "Image is not found", http.StatusNotFound)
		return
	}

	err = os.Remove(image.Path)
	if err != nil {
		log.Printf("image delete: %v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)

	http.Redirect(w, r, editPath, http.StatusFound)
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
	var data struct {
		Gallery GalleryOutput
		Images  []ImageOutput
	}
	gallery, err := g.galleryById(r, w, galleryBelongsToUser)
	if err != nil {
		return
	}

	data.Gallery = GalleryOutput{
		ID:    gallery.ID,
		Title: gallery.Title,
	}

	outputImages, err := g.outputImages(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	data.Images = outputImages

	g.Templates.Edit.Execute(w, r, data)
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

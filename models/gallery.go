package models

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB        *sql.DB
	ImagesDir string
}

func (service *GalleryService) Create(title string, userId int) (*Gallery, error) {
	gallery := Gallery{
		Title:  title,
		UserID: userId,
	}
	row := service.DB.QueryRow(`INSERT INTO galleries (user_id, title) VALUES ($1, $2) RETURNING id`, gallery.UserID, gallery.Title)
	err := row.Scan(&gallery.ID)
	return &gallery, err
}

func (service *GalleryService) ById(id int) (*Gallery, error) {
	gallery := Gallery{ID: id}
	row := service.DB.QueryRow(`SELECT user_id, title FROM galleries WHERE id = $1`, id)
	err := row.Scan(&gallery.UserID, &gallery.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, nil
	}
	return &gallery, err
}

func (service *GalleryService) ByUserId(userId int) ([]Gallery, error) {
	rows, err := service.DB.Query(`SELECT id, title FROM galleries WHERE user_id = $1`, userId)
	if err != nil {
		return nil, err
	}
	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{UserID: userId}
		err = rows.Scan(&gallery.ID, &gallery.Title)
		if err != nil {
			return nil, err
		}
		galleries = append(galleries, gallery)
	}
	err = rows.Err()
	if err != nil {
		return nil, rows.Err()
	}

	return galleries, err
}

func (service *GalleryService) Update(gallery *Gallery) error {
	_, err := service.DB.Exec(`
		UPDATE galleries 
		SET title = $2 
		WHERE id = $1`, gallery.ID, gallery.Title,
	)

	return err
}

func (service *GalleryService) Delete(gallery *Gallery) error {
	_, err := service.DB.Exec(`DELETE FROM galleries WHERE id = $1`, gallery.ID)
	return err
}

func (service *GalleryService) Images(galleryID int) ([]Image, error) {
	globPattern := filepath.Join(service.galleryDir(galleryID), "*")
	files, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("filepath.Glob error: %v", err)
	}
	var images []Image
	for _, file := range files {
		if service.hasExtension(file, service.extensions()) {
			images = append(images, Image{Path: file})
		}
	}
	return images, nil
}

func (service *GalleryService) galleryDir(id int) string {
	imagesDir := service.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}

func (service *GalleryService) extensions() []string {
	return []string{".jpg", ".png", ".jpeg", ".gif"}
}

func (service *GalleryService) hasExtension(path string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.ToLower(filepath.Ext(path)) == strings.ToLower(ext) {
			return true
		}
	}
	return false
}

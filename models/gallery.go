package models

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type Image struct {
	Path      string
	Filename  string
	Size      int64
	GalleryID int
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
	if err != nil {
		return fmt.Errorf("deleting from galleries: %v", err)
	}
	err = os.RemoveAll(service.galleryDir(gallery.ID))
	if err != nil {
		return fmt.Errorf("deleting gallery dir: %v", err)
	}
	return nil
}

func (service *GalleryService) Images(galleryID int) ([]Image, error) {
	globPattern := filepath.Join(service.galleryDir(galleryID), "*")
	files, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("filepath.Glob error: %v", err)
	}
	var images []Image
	for _, path := range files {
		if hasExtension(path, service.extensions()) {
			images = append(images, Image{
				Path:      path,
				Filename:  filepath.Base(path),
				GalleryID: galleryID,
			})
		}
	}
	return images, nil
}

func (service *GalleryService) Image(galleryID int, filename string) (Image, error) {
	imgPath := filepath.Join(service.galleryDir(galleryID), filename)
	fileInfo, err := os.Stat(imgPath)
	if err != nil {
		return Image{}, fmt.Errorf("searching an image: %v", err)
	}

	return Image{
		Path:      imgPath,
		Filename:  fileInfo.Name(),
		Size:      fileInfo.Size(),
		GalleryID: galleryID,
	}, nil
}

func (service *GalleryService) CreateImage(name string, galleryId int, file io.Reader) error {
	contentTypeSliceOfFile, err := checkContentType(file, service.imageContentTypes())
	if err != nil {
		return err
	}
	err = checkExtension(name, service.extensions())
	if err != nil {
		return err
	}
	dir := service.galleryDir(galleryId)
	imagePath := filepath.Join(dir, name)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("creating a gallery image directory: %v", err)
	}

	dst, err := os.Create(imagePath)
	defer dst.Close()
	if err != nil {
		return fmt.Errorf("creating a file to write an image to: %v", err)
	}

	completeFile := io.MultiReader(bytes.NewReader(contentTypeSliceOfFile), file)
	_, err = io.Copy(dst, completeFile)
	if err != nil {
		return fmt.Errorf("copying file contents to a disk: %v", err)
	}

	return nil
}

func (service *GalleryService) CreateImageFromUrl(galleryId int, url string) error {
	filename := path.Base(url)
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("file downloading error: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed downloading a file: status code %d", response.StatusCode)
	}

	defer response.Body.Close()

	return service.CreateImage(filename, galleryId, response.Body)
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

func (service *GalleryService) imageContentTypes() []string {
	return []string{"image/jpg", "image/jpeg", "image/png", "image/gif"}
}

func hasExtension(path string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.ToLower(filepath.Ext(path)) == strings.ToLower(ext) {
			return true
		}
	}
	return false
}

package models

import (
	"database/sql"
)

type Image struct {
	Path      string
	Filename  string
	Size      int64
	GalleryID int
}

type ImageService struct {
	DB *sql.DB
}

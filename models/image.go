package models

import (
	"database/sql"
)

type Image struct {
	Path string
}

type ImageService struct {
	DB *sql.DB
}

package main

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vityayka/go-zero/models"
)

func main() {
	gService := models.GalleryService{}
	images, err := gService.Images(2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("images: %v", images)
}

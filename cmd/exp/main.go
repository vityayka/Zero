package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// type User struct {
// 	Name string
// }

func main() {
	t, err := template.ParseFiles("cmd/exp/hello.gohtml")
	if err != nil {
		log.Printf("error parsing template %v", err)
		return
	}

	type Occupation struct {
		Salary   int32
		Position string
	}

	r := chi.NewRouter()
	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := struct {
			Name       string
			Occupation Occupation
		}{
			Name: "Ivan",
			Occupation: Occupation{
				Salary:   3000,
				Position: `<script>alert("azaz motherfucker")</script>`,
			},
		}

		err = t.Execute(w, user)
		if err != nil {
			panic(err)
		}
	}))

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)

}

package main

import (
	"fmt"
	"net/http"
	"time"
)

// func route(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/dashboard":
// 		handlerDash(w, r)
// 	case "/":
// 		defaultHandler(w)
// 	default:
// 		http.Error(w, "PAGE NOT FOUND", http.StatusNotFound)
// 		// w.WriteHeader(http.StatusNotFound)
// 		// fmt.Fprint(w, "<h1>Page not found!!!</h1>")
// 	}
// }

func defaultHandler(w http.ResponseWriter) {
	w.Header().Set("content-type", "text/html")
	fmt.Fprint(w, "<h1>Welcome12312312!!</h1>")
}

func handlerDash(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	w.Header().Set("content-type", "text/html")
	fmt.Fprint(w, "<pre>"+path+"</pre>")
}

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/dashboard":
		handlerDash(w, r)
	case "/":
		defaultHandler(w)
	default:
		http.Error(w, "PAGE NOT FOUND", http.StatusNotFound)
		// w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, "<h1>Page not found!!!</h1>")
	}
}

func main() {
	// mux := http.NewServeMux()
	// mux.HandleFunc()
	// http.HandleFunc("/", route)
	var router Router
	time.Sleep(1 * time.Second)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", router)
}

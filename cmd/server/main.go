package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/BTBurke/vatinator/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var (
	version string = "dev"
	commit  string
	date    string
)

func main() {
	// TODO: add viper config

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dir, err := ioutil.TempDir("", "vat-server")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Data directory: %s", dir)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Logger)
	r.Post("/file", handlers.FileAddHandler(dir))
	log.Printf("Serving running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

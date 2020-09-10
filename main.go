package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/drahoslavzan/mailtracker/database"

	"github.com/go-chi/chi"
)

const port = 8000
const prefix = "/seen/"
const suffix = ".png"
const preflen = len(prefix)
const pixel = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="

func getPixelImage() ([]byte, error) {
	return base64.StdEncoding.DecodeString(pixel)
}

func handleTrack(img []byte, emailRepo database.EmailRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSuffix(strings.ToLower(r.URL.Path[preflen:]), suffix)
		emailRepo.TrackSeen(id)
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		w.Write(img)
	}
}

func main() {
	router := chi.NewRouter()
	img, err := getPixelImage()
	if err != nil {
		log.Fatal(err)
	}

	dbCfg := database.GetConfig()
	db := database.NewDatabase(dbCfg)

	router.HandleFunc(prefix+"*", handleTrack(img, db.NewEmailRepository()))

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: router}

	log.Printf("serving HTTP on port: %d", port)
	log.Fatal(server.ListenAndServe())
}

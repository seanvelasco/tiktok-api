package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{user}", handleProfile)
	mux.HandleFunc("GET /{user}/video/{post}", handleVideo)
	mux.HandleFunc("GET /{user}/info/{post}", handlePost)
	mux.HandleFunc("GET /{user}/video/{post}/info", handlePost)
	mux.HandleFunc("GET /comments/{post}", handleComments)
	mux.HandleFunc("GET /{user}/comments/{post}", handleComments)
	mux.HandleFunc("GET /{user}/video/{post}/comments", handleComments)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}

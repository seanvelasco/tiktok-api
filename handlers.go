package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func handlePost(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user")
	postID := r.PathValue("post")

	post, err := getPost(userID, postID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleVideo(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("post")
	video, err := getVideo(postID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer video.Close()

	w.Header().Set("Content-Type", "video/mp4")

	_, err = io.Copy(w, video)

	if err != nil {
		http.Error(w, "Failed to write the video content", http.StatusInternalServerError)
		return
	}
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user")

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleComments(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("post")
	comments, err := getThread(postID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

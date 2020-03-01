package main

import (
	"log"
	"net/http"
)

// GetPostsHandler returns all posts
func (dbServer *Server) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	postsJSON, jsonErr := dbServer.api.GetPosts()

	if jsonErr != nil {
		log.Fatal(jsonErr)
		panic(jsonErr.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(postsJSON)
}

// GetPostAuthorHandler returns all authors
func (dbServer *Server) GetPostAuthorHandler(w http.ResponseWriter, r *http.Request) {
	authorsJSON, jsonErr := dbServer.api.GetAuthors()

	if jsonErr != nil {
		log.Fatal(jsonErr)
		panic(jsonErr.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(authorsJSON)
}

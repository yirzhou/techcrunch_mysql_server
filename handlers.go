package main

import (
	"encoding/json"
	"fmt"
	"log"
	"medium_mysql_server/api"
	"net/http"
)

// GetPostsHandler returns all posts
func (dbServer *Server) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "invalid_http_method")
		return
	}
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
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "invalid_http_method")
		return
	}
	authorsJSON, jsonErr := dbServer.api.GetAuthors()

	if jsonErr != nil {
		log.Fatal(jsonErr)
		panic(jsonErr.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(authorsJSON)
}

// CreatePost creates a new post
func (dbServer *Server) PostArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "invalid_http_method")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var reqPost api.Post
	if err := decoder.Decode(&reqPost); err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Print(reqPost)

	newPostId := dbServer.api.GetMaxPostId()
	fmt.Printf("New PostID: [%d]\n", newPostId)

}

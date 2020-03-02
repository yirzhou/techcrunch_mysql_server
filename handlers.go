package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"medium_mysql_server/api"
	"net/http"
	"time"
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

// CreatePost creates a new post, inserts tags, topics, etc.
func (dbServer *Server) PostArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "invalid_http_method")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var reqPost api.PostInfo
	if err := decoder.Decode(&reqPost); err != nil {
		log.Fatal(err)
		return
	}

	newPostId := dbServer.api.GetMaxPostId() + 1
	reqPost.Date = time.Now()
	reqPost.PostID = newPostId
	err := dbServer.api.InsertPost(reqPost)

	// Insert Tags, Topics, and Authors
	for _, tag := range reqPost.Tags {
		err = dbServer.api.InsertPostTag(api.Tag{PostID: newPostId, Tag: sql.NullString{String: tag, Valid: true}})
	}

	for _, topic := range reqPost.Topics {
		err = dbServer.api.InsertPostTopic(api.Topic{PostID: newPostId, Topic: sql.NullString{String: topic, Valid: true}})
	}

	for _, author := range reqPost.Authors {
		err = dbServer.api.InsertPostAuthor(api.Author{PostID: newPostId, AuthorID: author})
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Post with ID [%d] failed to be created\n", reqPost.PostID)
	} else {
		w.WriteHeader(http.StatusOK)
		log.Printf("Post with ID [%d] has been created\n", reqPost.PostID)
	}
}

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"medium_mysql_server/api"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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
		log.Println(jsonErr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(postsJSON)
}

// JoinGroupHandler adds a user to an existing group
func (dbServer *Server) JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "invalid_http_method")
		return
	}

	groupId, _ := strconv.ParseInt(mux.Vars(r)["groupId"], 10, 64)
	userId := r.FormValue("user_id")
	if err := dbServer.api.UpdateGroupWithUser(groupId, userId); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

// GetGroupsHandler returns information of available/existing groups.
func (dbServer *Server) GetGroupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "invalid_http_method")
		return
	}
	groupsJSON, jsonErr := dbServer.api.ListGroupsWithId()
	if jsonErr != nil {
		log.Println(jsonErr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(groupsJSON)
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

	var err error
	err = dbServer.api.InsertPost(reqPost)
	// Insert Category, Tags, Topics, and Authors
	if err := dbServer.api.InsertNewCategory(reqPost.Category); err != nil {
		log.Println(err.Error())
	}

	for _, tag := range reqPost.Tags {
		tagInfo := api.Tag{PostID: newPostId, Tag: sql.NullString{String: tag, Valid: true}}
		err = dbServer.api.InsertNewTag(tagInfo)
		err = dbServer.api.InsertPostTag(tagInfo)
	}

	for _, topic := range reqPost.Topics {
		topicInfo := api.Topic{PostID: newPostId, Topic: sql.NullString{String: topic, Valid: true}}
		err = dbServer.api.InsertNewTopic(topicInfo)
		err = dbServer.api.InsertPostTopic(topicInfo)
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

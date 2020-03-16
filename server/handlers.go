package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"medium_mysql_server/api"

	"github.com/gorilla/mux"
)

// CheckIfUserIsLoggedIn checks if the user is logged in.
func (server *Server) CheckIfUserIsLoggedIn(userId string, w *http.ResponseWriter) bool {
	if !server.api.IsUserLoggedIn(userId) {
		log.Println("User not logged in.")
		(*w).WriteHeader(http.StatusMethodNotAllowed)
		http.Error(*w, "User not logged in.", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

// GetPostsHandler returns all posts
func (server *Server) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodGet) {
		return
	}

	postsJSON, jsonErr := server.api.GetPosts()

	if jsonErr != nil {
		log.Println(jsonErr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(postsJSON)
}

// UserLogInHandler logs a user in.
func (server *Server) UserAuthHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodPost) {
		return
	}
	userId := mux.Vars(r)["userId"]
	action := mux.Vars(r)["action"]

	var err error
	if action == "login" {
		if err = server.api.AuthenticateUser(userId, action); err == nil {
			w.WriteHeader(http.StatusAccepted)
		}
	} else if action == "logout" {
		if err := server.api.AuthenticateUser(userId, action); err == nil {
			w.WriteHeader(http.StatusAccepted)
		}
	} else {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
	}

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
	}
}

// JoinGroupHandler adds a user to an existing group
func (server *Server) JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodPut) {
		return
	}

	groupId, _ := strconv.ParseInt(mux.Vars(r)["groupId"], 10, 64)
	userId := r.FormValue("user_id")
	if err := server.api.AddUserToGroup(groupId, userId); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

// UnfollowTopicHandler fulfils the user request of unfollowing a new topic.
func (server *Server) UnfollowTopicHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodPut) {
		return
	}

	userId := mux.Vars(r)["userId"]
	if !server.CheckIfUserIsLoggedIn(userId, &w) {
		return
	}

	topic := r.FormValue("topic")
	if err := server.api.RemoveTopic(userId, topic); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

// FollowTopicHandler fulfils the user request of following a new topic.
func (server *Server) FollowTopicHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodPut) {
		return
	}

	userId := mux.Vars(r)["userId"]
	if !server.CheckIfUserIsLoggedIn(userId, &w) {
		return
	}

	topic := r.FormValue("topic")
	if err := server.api.AddFollowedTopic(userId, topic); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

// ResponseToPostHandler handles a thumb-up/down from a user.
func (server *Server) ResponseToPostHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodPost) {
		return
	}

	postId := mux.Vars(r)["postId"]
	userId := r.FormValue("user_id")
	action := r.FormValue("action")

	log.Printf("postId: [%s], userId: [%s], action: [%s]\n", postId, userId, action)
	if !server.CheckIfUserIsLoggedIn(userId, &w) {
		return
	}

	if err := server.api.ThumbUpDownPost(postId, userId, action); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

// GetNewPostsForUserHandler retrieves the posts added after the time the user was logged in last time.
func (server *Server) GetNewPostsForUserHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodGet) {
		return
	}
	userId := mux.Vars(r)["userId"]
	if !server.CheckIfUserIsLoggedIn(userId, &w) {
		return
	}

	postsJSON, jsonErr := server.api.GetNewPostsForUser(userId)

	if jsonErr != nil {
		log.Println(jsonErr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(postsJSON)
}

// GetGroupsHandler returns information of available/existing groups.
func (server *Server) GetGroupsHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodGet) {
		return
	}

	groupsJSON, jsonErr := server.api.ListGroupsWithId()
	if jsonErr != nil {
		log.Println(jsonErr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(groupsJSON)
}

func (server *Server) GetFollowedTopicsHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodGet) {
		return
	}

	userId := mux.Vars(r)["userId"]
	topicsJSON, jsonErr := server.api.GetFollowedTopics(userId)
	if jsonErr != nil {
		log.Println(jsonErr.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(topicsJSON)
}

// GetPostAuthorHandler returns all authors
func (server *Server) GetPostAuthorHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodGet) {
		return
	}

	authorsJSON, jsonErr := server.api.GetAuthors()

	if jsonErr != nil {
		log.Println(jsonErr.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(authorsJSON)
}

// PostArticleHandler creates a new post, inserts tags, topics, etc.
func (server *Server) PostArticleHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodPost) {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var reqPost api.PostInfo
	if err := decoder.Decode(&reqPost); err != nil {
		log.Fatal(err)
		return
	}

	newPostId := server.api.GetMaxPostId() + 1

	reqPost.Date = time.Now()
	reqPost.PostID = newPostId

	var err error
	err = server.api.InsertPost(reqPost)
	// Insert Category, Tags, Topics, and Authors
	if err := server.api.InsertNewCategory(reqPost.Category); err != nil {
		log.Println(err.Error())
	}

	for _, tag := range reqPost.Tags {
		tagInfo := api.Tag{PostID: newPostId, Tag: sql.NullString{String: tag, Valid: true}}
		err = server.api.InsertNewTag(tagInfo)
		err = server.api.InsertPostTag(tagInfo)
	}

	for _, topic := range reqPost.Topics {
		topicInfo := api.Topic{PostID: newPostId, Topic: sql.NullString{String: topic, Valid: true}}
		err = server.api.InsertNewTopic(topicInfo)
		err = server.api.InsertPostTopic(topicInfo)
	}

	for _, author := range reqPost.Authors {

		err = server.api.InsertPostAuthor(api.Author{PostID: newPostId, AuthorID: author})
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Post with ID [%d] failed to be created\n", reqPost.PostID)
	} else {
		w.WriteHeader(http.StatusOK)
		log.Printf("Post with ID [%d] has been created\n", reqPost.PostID)
	}
}

func (server *Server) checkMethod(w *http.ResponseWriter, r *http.Request, correctMethod string) bool {
	if r.Method != correctMethod {
		(*w).WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(*w, "invalid_http_method")
		return false
	}
	return true
}

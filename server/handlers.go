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

// GetPostsHandler returns all posts.
func (server *Server) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodGet) {
		return
	}

	postsJSON, jsonErr := server.api.GetPosts()

	if jsonErr != nil {
		log.Println(jsonErr)
	}

	log.Println("posts fetched successfully.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(postsJSON)
}

// GetCategoriesHandler returns all categories.
func (server *Server) GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodGet) {
		return
	}

	categoriesJSON, jsonErr := server.api.GetCategories()
	if jsonErr != nil {
		log.Println(jsonErr)
	}
	log.Println("categories fetched successfully.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(categoriesJSON)

}

// GetTopicsHandler returns all categories.
func (server *Server) GetTopicsHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodGet) {
		return
	}

	topicsJSON, jsonErr := server.api.GetTopics()
	if jsonErr != nil {
		log.Println(jsonErr)
	}
	log.Println("topics fetched successfully.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(topicsJSON)

}

// UserLogInHandler logs a user in.
func (server *Server) UserAuthHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodPost) {
		return
	}

	userId := r.FormValue("user_id")
	password := r.FormValue("password")
	action := mux.Vars(r)["action"]

	var err error
	if action == "login" || action == "logout" {
		if err = server.api.AuthenticateUser(userId, password, action); err == nil {
			log.Printf("%s successful for [%s] \n", action, userId)
			w.WriteHeader(http.StatusAccepted)
		}
	} else {
		log.Println("wrong action.")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
	}

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprint(w, err.Error())
	}
}

// NewGroupHandler creates a new group and adds that user to the group.
func (server *Server) NewGroupHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodPost) {
		return
	}

	userId := r.FormValue("user_id")

	if err := server.api.CreateGroup(userId); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
	} else {
		log.Printf("adding user [%s] to new group successfully\n", userId)
		w.WriteHeader(http.StatusAccepted)
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
		log.Printf("join [%s] into [%d] successfully\n", userId, groupId)
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
		log.Printf("user [%s] unfollowing topic [%s] successfully\n", userId, topic)
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
		log.Printf("user [%s] following topic [%s] successfully\n", userId, topic)
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
		log.Printf("user [%s] giving post [%s] a [%s] successfully\n", userId, postId, action)
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

	log.Printf("successfully retrieved new posts for user [%s]\n", userId)
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

	log.Println("successfully retrieved groups")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(groupsJSON)
}

// GetThumbsHandler gets the number of thumbs of a post.
func (server *Server) GetThumbsHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodGet) {
		return
	}

	postId, err := strconv.ParseInt(mux.Vars(r)["postId"], 10, 64)
	if err != nil {
		log.Println(err)
	}

	thumbsJSON, jsonErr := server.api.GetPostThumbs(postId)
	if jsonErr != nil {
		log.Println(jsonErr)
	}

	log.Printf("successfully retrieved thumb count for post id [%d]\n", postId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(thumbsJSON)
}

// GetFollowedTopicsHandler returns the topics the user follows.
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

	userId := mux.Vars(r)["userId"]
	if !server.CheckIfUserIsLoggedIn(userId, &w) {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var reqPost api.PostInfo
	if err := decoder.Decode(&reqPost); err != nil {
		log.Println(err)
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
		if err = server.api.InsertPostAuthor(api.Author{PostID: newPostId, AuthorID: author}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("post with ID [%d] failed to be created\n", reqPost.PostID)
		} else {
			w.WriteHeader(http.StatusOK)
			log.Printf("post with ID [%d] has been created\n", reqPost.PostID)
		}
	}
}

// GetPostTopicsHandler handles getting PostTopics.
func (server *Server) GetPostTopicsHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodGet) {
		return
	}

	postTopicsJSON, jsonErr := server.api.GetPostTopics()

	if jsonErr != nil {
		log.Println(jsonErr)
	}

	log.Println("post topics fetched successfully.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(postTopicsJSON)
}

// UserSignUpHandler handles signning up users.
func (server *Server) UserSignUpHandler(w http.ResponseWriter, r *http.Request) {
	if !server.checkMethod(&w, r, http.MethodPost) {
		return
	}

	r.ParseForm()

	userId := r.FormValue("user_id")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	password := r.FormValue("password")

	if err := server.api.CreateUser(userId, firstName, lastName, password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("user [%s] failed to be created\n", userId)
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
	} else {
		w.WriteHeader(http.StatusAccepted)
		log.Printf("user [%s] has been created\n", userId)
	}

}

func (server *Server) checkMethod(w *http.ResponseWriter, r *http.Request, correctMethod string) bool {
	if r.Method != correctMethod {
		log.Printf("client method [%s] is invalid.\n", r.Method)
		(*w).WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(*w)
		return false
	}
	return true
}

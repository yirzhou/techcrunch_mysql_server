package server

import (
	"log"
	"medium_mysql_server/api"
	"net/http"

	"github.com/gorilla/mux"
)

// Server is a struct that will respond to various requests.
type Server struct {
	api    *api.API
	router *mux.Router
}

// RegisterServer creates a new server with an API to interact with database and a router.
func RegisterServer() *Server {
	server := &Server{api: api.NewAPI(), router: mux.NewRouter()}
	server.registerRoutes()
	return server
}

// Listen listens to 127.0.0.1:8080 for incoming requests.
func (server *Server) Listen() {
	log.Printf("Starting Listening on port [%d]...\n", 8080)
	http.ListenAndServe("127.0.0.1:8080", server.router)
}

func (server *Server) registerRoutes() {
	server.router.HandleFunc("/posts", server.GetPostsHandler)
	// GET: UserGroups
	server.router.HandleFunc("/groups", server.GetGroupsHandler)
	// GET: PostAuthors
	server.router.HandleFunc("/authors", server.GetPostAuthorHandler)
	// GET: FollowedTopics
	server.router.HandleFunc("/users/{userId}/topics", server.GetFollowedTopicsHandler)
	// POST: User Authentication
	server.router.HandleFunc("/users/{userId}/{action:(?:login|logout)}", server.UserAuthHandler)
	// POST: Post
	server.router.HandleFunc("/posts/new", server.PostArticleHandler)
	// PUT: Join UserGroup
	server.router.Path("/groups/{groupId:[0-9]+}/add").Queries("user_id", "{userId}").HandlerFunc(server.JoinGroupHandler)
	// PUT: Follow Topic
	server.router.Path("/users/{userId}/topics/add").Queries("topic", "{topic}").HandlerFunc(server.FollowTopicHandler)
}

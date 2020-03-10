package main

import (
	"log"
	"net/http"

	"medium_mysql_server/api"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Server is a struct that will respond to various requests.
type Server struct {
	api *api.API
}

func main() {

	dbServer := Server{api.NewAPI()}

	r := mux.NewRouter()
	// GET: Posts
	r.HandleFunc("/posts", dbServer.GetPostsHandler)
	// GET: UserGroups
	r.HandleFunc("/groups", dbServer.GetGroupsHandler)
	// GET: PostAuthors
	r.HandleFunc("/authors", dbServer.GetPostAuthorHandler)
	// POST: Post
	r.HandleFunc("/posts/new", dbServer.PostArticleHandler)
	// PUT: Join UserGroup
	r.Path("/groups/{groupId:[0-9]+}/add").Queries("user_id", "{userId}").HandlerFunc(dbServer.JoinGroupHandler)

	log.Printf("Starting Listening on port [%d]...\n", 8080)
	http.ListenAndServe("127.0.0.1:8080", r)
}

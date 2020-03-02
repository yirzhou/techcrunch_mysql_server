package main

import (
	"log"
	"net/http"

	"medium_mysql_server/api"

	_ "github.com/go-sql-driver/mysql"
)

// Server is a struct that will respond to various requests.
type Server struct {
	api *api.API
}

func main() {

	dbServer := Server{api.NewAPI()}

	mux := http.NewServeMux()
	mux.HandleFunc("/posts", dbServer.GetPostsHandler)
	mux.HandleFunc("/authors", dbServer.GetPostAuthorHandler)
	mux.HandleFunc("/posts/new", dbServer.PostArticleHandler)

	log.Printf("Starting Listening on port [%d]...\n", 8080)
	http.ListenAndServe("127.0.0.1:8080", mux)
}

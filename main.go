package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Server is a struct that will respond to various requests.
type Server struct {
	api *API
}

func main() {
	mux := http.NewServeMux()
	dbServer := Server{NewAPI()}
	mux.HandleFunc("/posts", dbServer.GetPostsHandler)
	http.ListenAndServe("127.0.0.1:8080", mux)

}

// GetPostsHandler returns all
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

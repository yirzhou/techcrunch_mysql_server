package main

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Post is a struct that contains the id and the content of a Post.
type Post struct {
	ID   int64
	Post string
}

// API is the API object to interact with MySQL server
type API struct {
	db *sql.DB
}

// NewAPI returns an API object that is used to interact with MySQL server
func NewAPI() *API {
	db, err := sql.Open("mysql", "go:zhou98@tcp(127.0.0.1:3307)/medium")
	if err != nil {
		panic(err.Error())
	}
	return &API{db}
}

// GetPosts retrieves all posts and display them.
func (api *API) GetPosts() ([]byte, error) {
	q := `select * from Posts;`

	rows, err := api.db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	posts := make([]*Post, 5)
	for rows.Next() {
		post := &Post{}
		if err := rows.Scan(&post.ID, &post.Post); err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
		log.Printf("ID: [%d], Post: [%s]\n", post.ID, post.Post)
	}

	jsonResponse, jsonError := json.Marshal(posts)
	return jsonResponse, jsonError
}

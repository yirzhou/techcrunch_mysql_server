package api

import (
	"database/sql"
	"time"
)

const DateFormat = "1998-03-30"

// Post is a struct that contains the id and the content of a Post.
type Post struct {
	PostID   int64          `json:"post_id"`
	Category sql.NullString `json:"category"`
	Content  string         `json:"content"`
	Date     time.Time      `json:"date"`
	ImageSrc sql.NullString `json:"image_src"`
	Section  sql.NullString `json:"section"`
	Title    string         `json:"title"`
	URL      sql.NullString `json:"url"`
}

// Author is an object of an author with postID and the authorID
type Author struct {
	PostID   int64  `json:"post_id"`
	AuthorID string `json:"author_id"`
}

// Tag is an object with postID and tag
type Tag struct {
	PostID int64  `json:"post_id"`
	Tag    string `json:"tag"`
}

// Topic is an object with postID and topic
type Topic struct {
	PostID int64  `json:"post_id"`
	Topic  string `json:"topic"`
}

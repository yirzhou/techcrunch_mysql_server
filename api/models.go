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
	PostID int64          `json:"post_id"`
	Tag    sql.NullString `json:"tag"`
}

// Topic is an object with postID and topic
type Topic struct {
	PostID int64          `json:"post_id"`
	Topic  sql.NullString `json:"topic"`
}

// PostInfo represents a Post with tags and topics.
type PostInfo struct {
	PostID   int64          `json:"post_id"`
	Authors  []string       `json:"authors"`
	Category sql.NullString `json:"category"`
	Content  string         `json:"content"`
	Date     time.Time      `json:"date"`
	ImageSrc sql.NullString `json:"image_src"`
	Section  sql.NullString `json:"section"`
	Title    string         `json:"title"`
	URL      sql.NullString `json:"url"`
	Tags     []string       `json:"tags"`
	Topics   []string       `json:"topics"`
}

// UserGroupInfo is the primitive data struct for retrieving information of available groups.
type UserGroupInfo struct {
	GroupID   int64  `json:"group_id"`
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UserGroupListing is the actual data struct returned from the server.
type UserGroupListing struct {
	GroupID int64  `json:"group_id"`
	Users   []User `json:"users"`
}

type User struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
}

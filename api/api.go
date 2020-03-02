package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	_ "github.com/go-sql-driver/mysql"
)

const InstanceId = "i-082a39858ce1d7279"
const Region = "us-east-1"
const Port = 3306
const User = "go:zhou98@tcp"

// API is the API object to interact with MySQL server
type API struct {
	db *sql.DB
}

// NewAPI returns an API object that is used to interact with MySQL server
func NewAPI() *API {
	connection := fmt.Sprintf("%s(%s:%d)/medium?parseTime=true", User, getServerIP(), Port)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err.Error())
	}
	return &API{db}
}

func (api *API) executeQuery(q string) *sql.Rows {
	rows, err := api.db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	return rows
}

// GetPosts retrieves all posts.
func (api *API) GetPosts() ([]byte, error) {
	q := `select * from Post limit 10;`
	rows := api.executeQuery(q)

	posts := make([]*Post, 5)
	for rows.Next() {
		post := &Post{}
		if err := rows.Scan(&post.PostID,
			&post.Category,
			&post.Content,
			&post.Date,
			&post.ImageSrc,
			&post.Section,
			&post.Title,
			&post.URL); err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}
	defer rows.Close()

	jsonResponse, jsonError := json.Marshal(posts)
	return jsonResponse, jsonError
}

// GetAuthors retrieves all posts and their corresponding authors.
func (api *API) GetAuthors() ([]byte, error) {
	q := `select * from PostAuthor limit 10;`
	rows := api.executeQuery(q)

	authors := make([]*Author, 5)
	for rows.Next() {
		author := &Author{}
		if err := rows.Scan(&author.PostID, &author.AuthorID); err != nil {
			log.Fatal(err)
		}
		authors = append(authors, author)
	}
	defer rows.Close()

	jsonResponse, jsonError := json.Marshal(authors)
	return jsonResponse, jsonError
}

// GetMaxPostId returns the largest postID found in DB.
// It can be used for creating a new post.
func (api *API) GetMaxPostId() int64 {
	q := `select max(postID) from Post;`
	rows := api.executeQuery(q)

	var maxPostId int64
	for rows.Next() {
		if err := rows.Scan(&maxPostId); err != nil {
			log.Fatal(err)
		}
	}
	defer rows.Close()
	return maxPostId
}

func getServerIP() string {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(Region),
	})

	if err != nil {
		log.Print(err.Error())
		panic(err)
	}

	svc := ec2.New(session)
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(InstanceId),
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	return *result.Reservations[0].Instances[0].PublicIpAddress
}

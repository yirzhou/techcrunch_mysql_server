package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	_ "github.com/go-sql-driver/mysql"
)

const InstanceID = "i-082a39858ce1d7279"

// API is the API object to interact with MySQL server
type API struct {
	db *sql.DB
}

// NewAPI returns an API object that is used to interact with MySQL server
func NewAPI() *API {
	db, err := sql.Open("mysql", "go:zhou98@tcp(54.173.164.111:3306)/medium?parseTime=true")
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

func getServerIP() {
	svc := ec2.New(session.New())
	input := &ec2.DescribeInstanceAttributeInput{
		Attribute:  aws.String("instanceType"),
		InstanceId: aws.String(InstanceID),
	}

	result, err := svc.DescribeInstanceAttribute(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

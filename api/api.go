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

const (
	InstanceId = "i-082a39858ce1d7279"
	Region     = "us-east-1"
	Port       = 3306
	User       = "go:zhou98@tcp"
)

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
		log.Println(err)
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
			log.Println(err)
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
			log.Println(err.Error())
		}
		authors = append(authors, author)
	}
	defer rows.Close()

	jsonResponse, jsonError := json.Marshal(authors)
	return jsonResponse, jsonError
}

// InsertPostTag inserts a new tag for a post to DB
func (api *API) InsertPostTag(tagInfo Tag) error {
	stmtIns, err := api.db.Prepare("insert into PostTag values( ?, ?)")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if _, err = stmtIns.Exec(
		tagInfo.PostID,
		tagInfo.Tag); err != nil {
		log.Println(err.Error())
	}

	defer stmtIns.Close()
	return err
}

// GetTopicCount returns the number of occurrences of a topic.
func (api *API) GetTopicCount(topicInfo Topic) int {
	q := fmt.Sprintf("select count(*) from Topic where topic='%s';", topicInfo.Topic.String)
	rows := api.executeQuery(q)
	topicCount := 0
	for rows.Next() {
		if err := rows.Scan(&topicCount); err != nil {
			log.Println(err)
		}
	}
	defer rows.Close()
	return topicCount
}

// CheckIfCategoryExists checks if the category already exists.
func (api *API) CheckIfCategoryExists(category sql.NullString) bool {
	q := fmt.Sprintf("select count(*) from Category where category='%s';", category.String)
	rows := api.executeQuery(q)
	categoryCount := 0
	for rows.Next() {
		if err := rows.Scan(&categoryCount); err != nil {
			log.Println(err)
		}
	}
	defer rows.Close()
	if categoryCount == 1 {
		return true
	}
	return false
}

// CheckIfTagExists checks if the tag already exists.
func (api *API) CheckIfTagExists(tagInfo Tag) bool {
	q := fmt.Sprintf("select count(*) from Tag where tag='%s';", tagInfo.Tag.String)
	rows := api.executeQuery(q)
	tagCount := 0
	for rows.Next() {
		if err := rows.Scan(&tagCount); err != nil {
			log.Println(err)
		}
	}
	defer rows.Close()
	if tagCount > 0 {
		return true
	}
	return false
}

// InsertNewTag creates a new tag in DB.
func (api *API) InsertNewTag(tagInfo Tag) error {
	if existed := api.CheckIfTagExists(tagInfo); !existed {
		stmtIns, err := api.db.Prepare("insert into Tag values( ? )")
		if err != nil {
			log.Println(err.Error())
			return err
		}

		if _, err = stmtIns.Exec(tagInfo.Tag.String); err != nil {
			log.Println(err.Error())
			return err
		}
		defer stmtIns.Close()
	}
	return nil
}

// InsertNewCategory creates a new category in DB.
func (api *API) InsertNewCategory(category sql.NullString) error {
	if existed := api.CheckIfCategoryExists(category); !existed {
		stmtIns, err := api.db.Prepare("insert into Category values( ? )")
		if err != nil {
			log.Println(err.Error())
			return err
		}

		if _, err = stmtIns.Exec(category.String); err != nil {
			log.Println(err.Error())
			return err
		}
		defer stmtIns.Close()
	}
	return nil
}

// InsertNewTopic creates a new topic in DB.
func (api *API) InsertNewTopic(topicInfo Topic) error {
	if topicCount := api.GetTopicCount(topicInfo); topicCount == 0 {
		stmtIns, err := api.db.Prepare("insert into Topic values( ? )")
		if err != nil {
			log.Println(err.Error())
			return err
		}

		if _, err = stmtIns.Exec(topicInfo.Topic.String); err != nil {
			log.Println(err.Error())
			return err
		}
		defer stmtIns.Close()
	}
	return nil
}

// InsertPostTopic inserts a new topic for a post to DB
func (api *API) InsertPostTopic(topicInfo Topic) error {
	stmtIns, err := api.db.Prepare("insert into PostTopic values( ?, ?)")
	if err != nil {
		log.Print(err.Error())
		return err
	}
	if _, err = stmtIns.Exec(
		topicInfo.PostID,
		topicInfo.Topic); err != nil {
		log.Print(err.Error())
	}

	defer stmtIns.Close()
	return err
}

// InsertPostAuthor inserts an author for a post to DB
func (api *API) InsertPostAuthor(author Author) error {
	stmtIns, err := api.db.Prepare("insert into PostAuthor values( ?, ?)")
	if err != nil {
		log.Print(err.Error())
		return err
	}
	if _, err = stmtIns.Exec(
		author.PostID,
		author.AuthorID); err != nil {
		log.Print(err.Error())
	}

	defer stmtIns.Close()
	return err
}

// InsertPost inserts a new post to DB.
func (api *API) InsertPost(postInfo PostInfo) error {
	stmtIns, err := api.db.Prepare("insert into Post values( ?, ?, ?, ?, ?, ?, ?, ? )")
	if err != nil {
		log.Print(err.Error())
		return err
	}
	_, err = stmtIns.Exec(
		postInfo.PostID,
		postInfo.Category,
		postInfo.Content,
		postInfo.Date,
		postInfo.ImageSrc,
		postInfo.Section,
		postInfo.Title,
		postInfo.URL)
	defer stmtIns.Close()
	return err
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

	serverIP := *result.Reservations[0].Instances[0].PublicIpAddress
	log.Printf("ServerIP Found: [%s]\n", serverIP)
	return serverIP
}

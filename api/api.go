package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	_ "github.com/go-sql-driver/mysql"
)

const (
	InstanceId = "i-082a39858ce1d7279"
	Region     = "us-east-1"
	Port       = 3306
	GoUser     = "go:zhou98@tcp"
	ThumbUp    = "thumbup"
	ThumbDown  = "thumbdown"
)

// API is the API object to interact with MySQL server
type API struct {
	db *sql.DB
}

// NewAPI returns an API object that is used to interact with MySQL server
func NewAPI() *API {
	connection := fmt.Sprintf("%s(%s:%d)/medium?parseTime=true", GoUser, getServerIP(), Port)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err.Error())
	}
	return &API{db}
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

func (api *API) executeQuery(q string) *sql.Rows {
	rows, err := api.db.Query(q)
	if err != nil {
		log.Println(err)
	}
	return rows
}

func (api *API) checkIfRecordExists(key, column, table string) bool {
	q := fmt.Sprintf("select count(1) from %s where %s=%s;", table, column, key)
	rows := api.executeQuery(q)
	existed := false

	for rows.Next() {
		if err := rows.Scan(&existed); err != nil {
			log.Println(err.Error())
		}
	}
	return existed
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

// AddUserToGroup will add a user to a group.
func (api *API) AddUserToGroup(groupId int64, userId string) error {
	q := fmt.Sprintf("insert into UserGroup values (%d, '%s');", groupId, userId)
	_, err := api.db.Query(q)
	return err
}

// AddFollowedTopic will add a topic to the followed topics of a user.
func (api *API) AddFollowedTopic(userId string, topic string) error {
	stmtIns, err := api.db.Prepare("insert into FollowTopic values( ?, ?)")
	if _, err := stmtIns.Exec(userId, topic); err != nil {
		log.Println(err)
	}
	return err
}

// RemoveTopic removes a topic for a user.
func (api *API) RemoveTopic(userId string, topic string) error {
	stmtIns, err := api.db.Prepare("delete from FollowTopic where userID=? and topic=?;")
	if _, err := stmtIns.Exec(userId, topic); err != nil {
		log.Println(err)
	}
	return err
}

// GetFollowedTopics retrieves all the topics the user follows.
func (api *API) GetFollowedTopics(userId string) ([]byte, error) {
	q := fmt.Sprintf(`select topic from FollowTopic where userID='%s';`, userId)
	rows := api.executeQuery(q)

	topics := make([]string, 0)
	var topic string
	for rows.Next() {
		if err := rows.Scan(&topic); err == nil {
			topics = append(topics, topic)
		} else {
			log.Println(err.Error())
		}
	}
	defer rows.Close()
	return json.Marshal(topics)
}

// GetNewPostsForUser retrieves all the topics the user follows.
func (api *API) GetNewPostsForUser(userId string) ([]byte, error) {

	/*
		select * from Post inner join PostTopic using (postID) inner join FollowTopic using (topic), (select lastLoggedIn from User where userID='yirzhou') as T where userID='yirzhou' and T.lastLoggedIn <= Post.date;
	*/

	q := fmt.Sprintf(`select postID, category, content, date, img_src, section, title, url
		from Post inner join PostTopic using (postID) inner join FollowTopic using (topic), 
		(select lastLoggedIn from User where userID='%s') as T where userID='%s' and T.lastLoggedIn <= Post.date;`,
		userId,
		userId)

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
	return json.Marshal(posts)

}

// ListGroupsWithId returns information of available groups.
func (api *API) ListGroupsWithId() ([]byte, error) {
	q := `select groupID, userID, firstName, lastName from UserGroup inner join User using (userID) order by groupID asc;`
	rows := api.executeQuery(q)

	groups := make(map[int64]*UserGroupListing)
	for rows.Next() {
		var groupID int64
		var userID, firstName, lastName string

		if err := rows.Scan(&groupID, &userID, &firstName, &lastName); err == nil {
			if groupListing, ok := groups[groupID]; !ok {
				group := &UserGroupListing{GroupID: groupID, Users: make([]User, 0)}
				group.Users = append(group.Users, User{UserID: userID, FullName: fmt.Sprintf("%s %s", firstName, lastName)})
				groups[groupID] = group
			} else {
				groupListing.Users = append(groupListing.Users, User{UserID: userID, FullName: fmt.Sprintf("%s %s", firstName, lastName)})
			}
		} else {
			log.Println(err.Error())
		}
	}
	defer rows.Close()
	return json.Marshal(groups)
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

// IsUserLoggedIn returns a boolean value of the user authenticatiion status.
func (api *API) IsUserLoggedIn(userId string) bool {
	q := fmt.Sprintf("select isLoggedIn from User where userID='%s';", userId)
	rows := api.executeQuery(q)
	var isLoggedIn bool
	for rows.Next() {
		if err := rows.Scan(&isLoggedIn); err != nil {
			return isLoggedIn
		}
	}
	return isLoggedIn
}

// AuthenticateUser logs a user in/out depending on the action.
func (api *API) AuthenticateUser(userId string, action string) error {
	var q string
	if action == "login" {
		q = "update User set lastLoggedIn=?, isLoggedIn=1 where userID=?;"
	} else {
		q = "update User set lastLoggedIn=?, isLoggedIn=0 where userID=?;"
	}

	stmtIns, err := api.db.Prepare(q)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if _, err = stmtIns.Exec(time.Now(), userId); err != nil {
		log.Println(err)
	}
	defer stmtIns.Close()
	return err
}

// InsertPostTag inserts a new tag for a post to DB
func (api *API) InsertPostTag(tagInfo Tag) error {
	stmtIns, err := api.db.Prepare("insert into PostTag values( ?, ?)")

	if _, err = stmtIns.Exec(tagInfo.PostID, tagInfo.Tag); err != nil {
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
			log.Println(err)
		}
	}
	defer rows.Close()
	return maxPostId
}

// GetMaxGroupId returns the max groupID from UserGroup.
// It is used to create a new group
func (api *API) GetMaxGroupId() int64 {
	q := `select max(groupID) from UserGroup;`
	rows := api.executeQuery(q)

	var maxGroupId int64
	for rows.Next() {
		if err := rows.Scan(&maxGroupId); err != nil {
			log.Println(err.Error())
		}
	}
	defer rows.Close()
	return maxGroupId
}

// ThumbUpDownPost adds a thumb-up or down to a Post.
func (api *API) ThumbUpDownPost(postId, userId, action string) error {
	var qUser, qPost string

	if !api.checkIfRecordExists(postId, "postID", "Post") {
		return errors.New("invalid request")
	}

	if action == ThumbUp {
		qUser = `insert into UserThumbUp(userID, postID) values (?,?);`
		if api.checkIfRecordExists(postId, "postID", "PostThumbUp") {
			qPost = `update PostThumbUp set upCount=upCount+1 where postID=?;`
		} else {
			qPost = `insert into PostThumbUp(postID, upCount) values (?, 1);`
		}

		// Remove the thumbdown record for this (userID, postID)
		stmtRemove, _ := api.db.Prepare("delete from UserThumbDown where userID=? and postID=?;")
		_, _ = stmtRemove.Exec(userId, postId)
		defer stmtRemove.Close()

	} else if action == ThumbDown {
		qUser = `insert into UserThumbDown(userID, postID) values (?,?);`
		if api.checkIfRecordExists(postId, "postID", "PostThumbUp") {
			qPost = `update PostThumbUp set upCount=upCount-1 where postID=?;`
		} else {
			qPost = `insert into PostThumbUp(postID, upCount) values (?, -1);`
		}

		// Remove the thumbup record for this (userID, postID)
		stmtRemove, _ := api.db.Prepare("delete from UserThumbUp where userID=? and postID=?;")
		_, _ = stmtRemove.Exec(userId, postId)
		defer stmtRemove.Close()
	} else {
		return errors.New("invalid request")
	}

	stmtInsUser, err := api.db.Prepare(qUser)
	_, err = stmtInsUser.Exec(userId, postId)
	defer stmtInsUser.Close()

	if err != nil {
		return err
	}

	stmtInsPost, err := api.db.Prepare(qPost)
	_, err = stmtInsPost.Exec(postId)
	defer stmtInsPost.Close()
	return err
}

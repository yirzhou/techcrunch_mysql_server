package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

// CreateUser creates a new user.
func (api *API) CreateUser(userId, firstName, lastName, password string) error {
	q := `insert into User(userID, firstName, lastName, password) values (?, ?, ?, ?)`
	stmtIns, err := api.db.Prepare(q)

	if len(password) == 0 {
		if _, err := stmtIns.Exec(userId, firstName, lastName, nil); err != nil {
			log.Println(err)
		}
	} else {
		if _, err := stmtIns.Exec(userId, firstName, lastName, password); err != nil {
			log.Println(err)
		}
	}

	return err
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
func (api *API) AuthenticateUser(userId, password, action string) error {
	var q string
	if action == "login" {
		if storedPass := api.GetColumnFromTable("password", "User", "userID", userId); storedPass != nil && storedPass.(string) != password {
			return errors.New("your password/username is wrong")
		}

		q = "update User set lastLoggedIn=?, isLoggedIn=1 where userID=?"
	} else {
		q = "update User set lastLoggedIn=?, isLoggedIn=0 where userID=?;"
	}

	stmtIns, err := api.db.Prepare(q)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	var sqlResult sql.Result

	if sqlResult, err = stmtIns.Exec(time.Now(), userId); err != nil {
		log.Println(err)
	}

	if affectedRows, _ := sqlResult.RowsAffected(); affectedRows == 0 {
		err = fmt.Errorf("user id [%s] does not exist", userId)
	}

	defer stmtIns.Close()

	return err
}

// CreateGroup will create a group and add the user to this group.
func (api *API) CreateGroup(userId string) error {
	q := `insert into UserGroup (userID) values (?);`
	stmtIns, err := api.db.Prepare(q)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	var sqlResult sql.Result

	if sqlResult, err = stmtIns.Exec(userId); err != nil {
		log.Println(err)
	}

	if affectedRows, _ := sqlResult.RowsAffected(); affectedRows == 0 {
		err = fmt.Errorf("failed to add user [%s] to the new group", userId)
	}

	defer stmtIns.Close()

	return err
}

// AddUserToGroup will add a user to a group.
func (api *API) AddUserToGroup(groupId int64, userId string) error {
	q := fmt.Sprintf("insert into UserGroup values (%d, '%s');", groupId, userId)
	_, err := api.db.Query(q)
	return err
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

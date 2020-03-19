package api

import (
	"database/sql"
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

func (api *API) GetColumnFromTable(column, table, key, value string) interface{} {
	q := fmt.Sprintf("select %s from %s where %s='%s';", column, table, key, value)
	rows := api.executeQuery(q)

	var attrVal interface{}

	for rows.Next() {
		if err := rows.Scan(&attrVal); err != nil {
			log.Println(err.Error())
		}
	}
	return attrVal
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"medium_mysql_server/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	_ "github.com/go-sql-driver/mysql"
)

// Server is a struct that will respond to various requests.
type Server struct {
	api *api.API
}

func main() {
	getServerIP()
	dbServer := Server{api.NewAPI()}

	mux := http.NewServeMux()
	mux.HandleFunc("/posts", dbServer.GetPostsHandler)
	mux.HandleFunc("/authors", dbServer.GetPostAuthorHandler)

	log.Printf("Starting Listening on port [%d]...\n", 8080)
	http.ListenAndServe("127.0.0.1:8080", mux)
}

func getServerIP() {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		log.Print(err.Error())
		panic(err)
	}

	svc := ec2.New(session)
	input := &ec2.DescribeInstanceAttributeInput{
		Attribute:  aws.String("instanceType"),
		InstanceId: aws.String(api.InstanceID),
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

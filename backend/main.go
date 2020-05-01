package main

import (
	"medium_mysql_server/server"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	server := server.RegisterServer()
	server.Listen()
}

package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type Message struct {
	Email   string `db:"email" json:"email"`
	Title   string `db:"title" json:"title"`
	Content string `db:"content" json:"content"`
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func main() {
	r := gin.Default()
	v1 := r.Group("api/v1")

	{
		v1.POST("/message", PostMessage)
		v1.POST("/send", PostSend)
		v1.GET("/messages/:email", GetMessages)
	}

	r.Run(":8080")
}

func GetCassandraIP() (IP string) {

	IP = "127.0.0.1"
	fmt.Println(" Service cassandra IP : ", IP)
	return

}

func GetMessages(c *gin.Context) {

	// connect to the cluster
	cluster := gocql.NewCluster(GetCassandraIP())

	// A keyspace in Cassandra is a namespace that defines data replication on nodes. A cluster contains one keyspace per node.
	cluster.Keyspace = "demo"

	// May use gocql.Quorum
	cluster.Consistency = gocql.LocalOne
	session, _ := cluster.CreateSession()

	// Make sure that the connection can close once you are done.
	defer session.Close()

	//    ####################################### Query Logic ##########################################

	email := c.Params.ByName("email")

	fmt.Println("email : ", email)

	var messages []Message
	var message Message

	iter := session.Query("SELECT email FROM message WHERE email=?", email).Iter()
	for iter.Scan(&message.Email) {

		// fmt.Println("1 : ", message.Class, message.Name, message.Blob)
		messages = append(messages, message)
	}

	if len(messages) > 0 {
		c.JSON(200, messages)
	} else {
		c.JSON(404, gin.H{"error": "no message(s) into the table"})
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	// curl -i http://localhost:8080/api/messages/common
}

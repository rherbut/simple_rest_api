package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type Message struct {
	Email       string `db:"email" json:"email"`
	Title       string `db:"title" json:"title"`
	Content     string `db:"content" json:"content"`
	Magicnumber string `db:"magicnumber" json:"magicnumber"`
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func main() {
	r := gin.Default()
	v1 := r.Group("api")

	{
		v1.POST("/message", PostMessage)
		v1.POST("/send", PostSend)
		v1.GET("/messages/:email", GetMessages)
	}

	r.Run(":8080")
}

func GetCassandraIP() (IP string) {

	//var IP string
	/*
	   //Get a new client
	   client, err := api.NewClient(api.DefaultConfig())
	   if err != nil {
	      panic(err)
	   }

	    agent := client.Agent()
	    services, err := agent.Services()

	    fmt.Println("services : ", services)
	        if err != nil {
	                log.Fatal("err: %v", err)
	        }

	        if _, ok := services["cassandra"]; !ok {
	                log.Fatal("missing service: %v", services)
	        }

	    IP = services["cassandra"].Address
	    fmt.Println(" Service cassandra IP : ", services["cassandra"].Address)
	*/
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

	iter := session.Query("SELECT email, title, content, magicnumber FROM message WHERE email=?", email).Iter()
	for iter.Scan(&message.Email, &message.Title, &message.Content, &message.Magicnumber) {

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

package main

import (
	"fmt"

	gin "github.com/gin-gonic/gin"
	leveldb "github.com/syndtr/goleveldb/leveldb"
)

type Message struct {
	Message string
}

var db *leveldb.DB

func main() {
	fmt.Println("MiniQ is running...")

	var err error
	db, err = leveldb.OpenFile("./testDB", nil)
	if err != nil {
		fmt.Printf("there was an error")
		panic("panicking")
	}
	defer db.Close()

	r := gin.Default()

	//Returns the next message in the Queue
	r.GET("/getNextMessage/:key", getMessageHandler())
	//Adds a new message to the queue
	r.POST("/postMessage", AddMessageHandler())

	r.Run(":8080")
}

func Get(key string) string {
	data, _ := db.Get([]byte(key), nil)

	return string(data)
}

func Add(key []byte, val []byte) bool {
	err := db.Put(key, val, nil)

	if err != nil {
		panic(err)
	}

	return true
}

func getMessageHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		key := c.Param("key")
		value := Get(key)

		c.JSON(200, gin.H{
			"message": value,
		})
	}
	return gin.HandlerFunc(fn)
}

func AddMessageHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var newMessage Message
		if err := c.BindJSON(&newMessage); err != nil {
			fmt.Println("There was an error binding json")
			return
		}
		key := GenerateKey()

		success := Add([]byte(key), []byte(newMessage.Message))

		if !success {
			return
		}
	}

	return gin.HandlerFunc(fn)
}

func GenerateKey() string {
	return "xyz"
}

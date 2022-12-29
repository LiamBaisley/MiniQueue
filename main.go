package main

import (
	"flag"
	"fmt"
	"net/http"

	gin "github.com/gin-gonic/gin"
	leveldb "github.com/syndtr/goleveldb/leveldb"
)

type Message struct {
	Message string
}

type Key struct {
	Key string
}

var db *leveldb.DB
var config Config

// Keys are generated as unique and incremental so that we can leverage the fact that LevelDB stores key value pairs
// in order based on the key. using characters "a"-"z" we have of 26^15 or 1,677,259,342,285,725,925,376 possible keys.
// Keys also reset if the queue is emptied. Based on this we can assume that we should never run out of keys.
var firstKey = "aaaaaaaaaaaaaaa"

const ConfigFileName = "config.json"

func main() {
	var secret string
	var existing bool
	var hashErr error
	fmt.Println("MiniQ is running...")

	flag.StringVar(&secret, "s", "secret", "Used to configure the secret used to verify communications with the Queue.")
	flag.BoolVar(&existing, "e", false, "Whether or not to use an existing configuration.")

	flag.Parse()

	if existing && CheckFileExist(ConfigFileName) {
		fmt.Println("Found existing config file. Using existing config.")
		config = GetConfig(ConfigFileName)
	} else if !existing && secret != "" || secret != "secret" {
		config.SecurityHash, hashErr = CreateHash(secret)

		if hashErr != nil {
			panic("Could not hash security string.")
		}

		if WriteConfig(config, ConfigFileName) {
			panic("Could not write config file. Stopping program")
		}
	} else {
		panic("You need to configure a secret.")
	}

	var err error
	db, err = leveldb.OpenFile("./testDB", nil)
	if err != nil {
		fmt.Printf("there was an error")
		panic("panicking")
	}
	defer db.Close()

	r := gin.Default()

	r.Use(messageSizeMiddleware)
	r.Use(authMiddleware)

	//Returns the next message in the Queue
	r.GET("/message", getMessageHandler)
	//Adds a new message to the queue
	r.POST("/message", AddMessageHandler)
	//Confirm consumption of a message before deletion
	r.POST("/confirm", ConfirmConsumptionHandler)

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

func Delete(key []byte) bool {
	err := db.Delete(key, nil)

	if err != nil {
		panic(err)
	}

	return true
}

// Message size limitation to 5mb
func messageSizeMiddleware(c *gin.Context) {
	size := c.Request.ContentLength

	if size > 5000000 {
		c.AbortWithStatusJSON(413, gin.H{"error": "Message too large"})
	}

	c.Next()
}

// Authorize the user
func authMiddleware(c *gin.Context) {
	var secret = c.Request.Header["Authorization"]
	if result, _ := CompareHash(config.SecurityHash, secret[0]); !result {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unathorized"})
	}
	c.Next()
}

// Handler for getting the next message in the queue
func getMessageHandler(c *gin.Context) {
	iter := db.NewIterator(nil, nil)
	iter.First()
	Value := iter.Value()

	if Value == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "message not found"})
	} else {

		c.JSON(http.StatusOK, gin.H{
			"message": string(Value),
			"key":     string(iter.Key()),
		})
	}
}

// Handler for adding messages
func AddMessageHandler(c *gin.Context) {
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

	c.IndentedJSON(http.StatusOK, key)
}

func ConfirmConsumptionHandler(c *gin.Context) {
	var Key Key
	if err := c.BindJSON(&Key); err != nil {
		fmt.Println("There was an error binding json")
		return
	}

	Delete([]byte(Key.Key))

	c.IndentedJSON(http.StatusOK, Key)
}

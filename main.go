package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	gin "github.com/gin-gonic/gin"
	leveldb "github.com/syndtr/goleveldb/leveldb"
)

type Message struct {
	Message   string
	Timestamp string
}

type Key struct {
	Key string
}

var db *leveldb.DB
var secret string
var confirmConsume bool

// Keys are generated as unique and incremental so that we can leverage the fact that LevelDB stores key value pairs
// in order based on the key. The characters of the string are appended to the current date represented as a number in this format: 202319-aaaaaaaaaaaaaaa.
// This means that for any given day we could have 26^15 keys and should never have to worry about running out of keys.
var firstKey = "aaaaaaaaaaaaaaa"

const ConfigFileName = "config.json"

func main() {
	fmt.Println("MiniQ is running...")
	flag.BoolVar(&confirmConsume, "c", true, "Whether or not MiniQ should expect users of the Queue to confirm consumption of messages")
	var secretErr error

	secret, secretErr = GetEnvSecret()

	if secretErr != nil {
		panic("No environment secret set.")
	}

	var err error
	db, err = leveldb.OpenFile("./testDB", nil)
	if err != nil {
		fmt.Printf("there was an error opening the DB file %s", err)
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

func Get(key string) (string, error) {
	data, err := db.Get([]byte(key), nil)
	if err != nil {
		return "", fmt.Errorf("GetKey: there was an error retrieving the key %w", err)
	}
	return string(data), nil
}

func Add(key []byte, val []byte) error {
	err := db.Put(key, val, nil)

	if err != nil {
		return fmt.Errorf("AddKey: There was an error adding the key %w", err)
	}

	if !confirmConsume {
		err := db.Delete(key, nil)
		if err != nil {
			return fmt.Errorf("AddKey: there was an error removing the element from the queue %w", err)
		}
	}

	return nil
}

func Delete(key []byte) error {
	err := db.Delete(key, nil)

	if err != nil {
		return fmt.Errorf("DeleteKey: There was an error removing the key %w", err)
	}

	return nil
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
	var authSecret = c.Request.Header["Authorization"]
	if authSecret[0] != secret {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
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
		message := Message{}
		err := json.Unmarshal(Value, &message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "there was a server error"})
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   string(message.Message),
			"timestamp": string(message.Timestamp),
			"key":       string(iter.Key()),
		})
	}
}

// Handler for adding messages
func AddMessageHandler(c *gin.Context) {
	var newMessage Message
	if err := c.BindJSON(&newMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}
	key, err := GenerateKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}
	newMessage.Timestamp = time.Now().Format(time.RFC3339Nano)
	content, err := json.Marshal(newMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}
	adderr := Add([]byte(key), []byte(content))
	if adderr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": adderr})
	}

	c.IndentedJSON(http.StatusOK, key)
}

func ConfirmConsumptionHandler(c *gin.Context) {
	var Key Key
	if err := c.BindJSON(&Key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	err := Delete([]byte(Key.Key))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	c.IndentedJSON(http.StatusOK, Key)
}

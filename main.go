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

var db *leveldb.DB
var config Config

// Keys are generated as unique and incremental so that we can leverage the fact that LevelDB stores key value pairs
// in order based on the key. using characters "a"-"z" we have of 26^15 or 1,677,259,342,285,725,925,376 possible keys.
// Keys also reset if the queue is emptied. Based on this we can assume that we should never run out of keys.
var firstKey = "aaaaaaaaaaaaaaa"

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
		config = GetConfig()
	} else if !existing && secret != "" || secret != "secret" {
		config.SecurityHash, hashErr = CreateHash(secret)

		if hashErr != nil {
			panic("Could not hash security string.")
		}

		if WriteConfig(config) {
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

	//Returns the next message in the Queue
	r.GET("/getNextMessage", getMessageHandler)
	//Adds a new message to the queue
	r.POST("/postMessage", AddMessageHandler)
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

// Handler for getting the next message in the queue
func getMessageHandler(c *gin.Context) {
	var secret = c.Request.Header["Authorization"]
	if result, _ := CompareHash(config.SecurityHash, secret[0]); result {
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
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized.",
		})
	}
}

// Handler for adding messages
func AddMessageHandler(c *gin.Context) {
	var newMessage Message
	var secret = c.Request.Header["Authorization"]
	if result, _ := CompareHash(config.SecurityHash, secret[0]); result {
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
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized.",
		})
	}
}

func ConfirmConsumptionHandler(c *gin.Context) {
	var secret = c.Request.Header["Authorization"]
	var key string
	if result, _ := CompareHash(config.SecurityHash, secret[0]); result {
		if err := c.BindJSON(&key); err != nil {
			fmt.Println("There was an error binding json")
			return
		}

		Delete([]byte(key))

		c.IndentedJSON(http.StatusOK, key)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized.",
		})
	}
}

// Uses the first key if there are no messages in the queue, otherwise incrememts and returns the last key in the queue.
func GenerateKey() string {
	iter := db.NewIterator(nil, nil)
	first := iter.First()
	var newKey string

	if !first {
		return firstKey
	} else {
		iter.Last()
		currKey := iter.Key()
		newKey = IncrementKey(currKey, currKey[len(currKey)-1], len(currKey)-1)
	}

	return newKey
}

// Recursively goes through the characters in the key to determine which to increment.
func IncrementKey(key []byte, currentByte byte, index int) string {
	if index == 0 && int(currentByte) == 122 {
		//Considering we have such a large number of possible keys, we should never get here.
		//But we handle it just in case.
		panic("Out of keys exception")
	}

	if int(currentByte) == 122 {
		key[index] = byte(97)
		IncrementKey(key, key[index-1], index-1)
	} else {
		key[index]++
	}

	return string(key)
}

package main

import (
	"fmt"
	"net/http"

	gin "github.com/gin-gonic/gin"
	leveldb "github.com/syndtr/goleveldb/leveldb"
)

type Message struct {
	Message string
}

var db *leveldb.DB

// Keys are generated as unique and incremental so that we can leverage the fact that LevelDB stores key value pairs
// in order based on the key. using characters "a"-"z" we have of 26^15 or 1,677,259,342,285,725,925,376 possible keys.
// Keys also reset if the queue is emptied. Based on this we can assume that we should never run out of keys.
var firstKey = "aaaaaaaaaaaaaaa"

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
	r.GET("/getNextMessage", getMessageHandler)
	//Adds a new message to the queue
	r.POST("/postMessage", AddMessageHandler)

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
	iter := db.NewIterator(nil, nil)
	iter.First()
	Value := iter.Value()
	defer Delete(iter.Key())
	fmt.Printf("the value is %v", string(Value))
	c.JSON(http.StatusOK, gin.H{
		"message": string(Value),
	})
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

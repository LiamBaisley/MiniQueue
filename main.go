package main

import (
	"fmt"

	leveldb "github.com/syndtr/goleveldb/leveldb"
)

var db *leveldb.DB

func main() {
	fmt.Println("MiniQ is running...")
	db = GetConnection()
	defer db.Close()
	var testKey string = "hello"
	var testValue string = "world"

	fmt.Println("Adding a record")
	Add([]byte(testKey), []byte(testValue))

	fmt.Println("trying to retrieve record")
	var data = Get(testKey)

	fmt.Println(string(data))
}

func Get(key string) []byte {
	data, err := db.Get([]byte(key), nil)
	if err != nil {
		panic(err)
	}

	return data
}

func Add(key []byte, val []byte) bool {
	err := db.Put([]byte("key"), []byte("value"), nil)

	if err != nil {
		panic(err)
	}

	return true
}

func GetConnection() *leveldb.DB {
	db, _ := leveldb.OpenFile("./", nil)
	return db
}

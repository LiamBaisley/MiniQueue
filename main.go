package main

import (
	"fmt"

	leveldb "github.com/syndtr/goleveldb/leveldb"
)

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

	var testKey string = "hello"
	var testValue string = "world"

	fmt.Println("Adding a record")
	Add([]byte(testKey), []byte(testValue))

	fmt.Println("trying to retrieve record")
	var data = Get(testKey)

	fmt.Println(string(data))
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

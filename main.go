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
	Add([]byte("D"), []byte(testValue))
	Add([]byte("B"), []byte(testValue))
	Add([]byte("F"), []byte(testValue))
	Add([]byte("Z"), []byte(testValue))
	Add([]byte("1"), []byte(testValue))
	Add([]byte("7"), []byte(testValue))
	Add([]byte("A"), []byte(testValue))

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		fmt.Printf(string(key) + " " + string(value) + "\n")
	}
	iter.Release()

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

package main

import (
	"fmt"
	"strings"
	"time"
)

// Uses the first key if there are no messages in the queue, otherwise incrememts and returns the last key in the queue.
func GenerateKey() string {
	iter := db.NewIterator(nil, nil)
	first := iter.First()
	var newKey string

	if !first {
		return fmt.Sprintf("%v-%v", getCurrentDateAsString(), firstKey)
	} else {
		iter.Last()
		currKey := iter.Key()
		splitKey := strings.Split(string(currKey), "-")
		keyDate := splitKey[0]
		if keyDate == getCurrentDateAsString() {
			keyToIncrement := []byte(splitKey[1])
			newKey = IncrementKey(keyToIncrement, keyToIncrement[len(keyToIncrement)-1], len(keyToIncrement)-1)
			return fmt.Sprintf("%v-%v", getCurrentDateAsString(), newKey)
		} else {
			return fmt.Sprintf("%v-%v", getCurrentDateAsString(), firstKey)
		}
	}
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

func getCurrentDateAsString() string {
	year, month, day := time.Now().Date()
	return fmt.Sprintf("%v%v%v", year, int(month), day)
}

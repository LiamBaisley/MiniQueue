package main

import (
	"fmt"
	"strings"
	"time"
)

// Uses the first key if there are no messages in the queue, otherwise incrememts and returns the last key in the queue.
func GenerateKey() (string, error) {
	iter := db.NewIterator(nil, nil)
	first := iter.First()

	if !first {
		return fmt.Sprintf("%v-%v", getCurrentDateAsString(), firstKey), nil
	} else {
		iter.Last()
		currKey := iter.Key()
		splitKey := strings.Split(string(currKey), "-")
		keyDate := splitKey[0]
		if keyDate == getCurrentDateAsString() {
			keyToIncrement := []byte(splitKey[1])
			newKey, err := IncrementKey(keyToIncrement, keyToIncrement[len(keyToIncrement)-1], len(keyToIncrement)-1)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%v-%v", getCurrentDateAsString(), newKey), nil
		} else {
			return fmt.Sprintf("%v-%v", getCurrentDateAsString(), firstKey), nil
		}
	}
}

// Recursively goes through the characters in the key to determine which to increment.
func IncrementKey(key []byte, currentByte byte, index int) (string, error) {
	if index == 0 && int(currentByte) == 122 {
		//Considering we have such a large number of possible keys, we should never get here.
		//But we handle it just in case.
		return "", fmt.Errorf("IncrementKey: Out of keys exception")
	}

	if int(currentByte) == 122 {
		key[index] = byte(97)
		IncrementKey(key, key[index-1], index-1)
	} else {
		key[index]++
	}

	return string(key), nil
}

func getCurrentDateAsString() string {
	year, month, day := time.Now().Date()
	return fmt.Sprintf("%v%v%v", year, int(month), day)
}

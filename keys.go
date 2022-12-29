package main

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

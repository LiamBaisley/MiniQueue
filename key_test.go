package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrementKey_IfLastKey_ShouldReturnError(t *testing.T) {
	key := "zzzzzzzzzzzzzzz"
	_, err := IncrementKey([]byte(key), key[len(key)-1], 0)
	assert.Error(t, err)
}

func TestIncrementKey_GivenKey_ShouldIncrement(t *testing.T) {
	key := "aaaaaaaaaaaaaaa"

	updatedKey, _ := IncrementKey([]byte(key), key[len(key)-1], len(key)-1)
	assert.Equal(t, "aaaaaaaaaaaaaab", updatedKey)
}

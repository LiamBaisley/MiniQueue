package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrementKey_IfLastKey_ShouldPanic(t *testing.T) {
	key := "zzzzzzzzzzzzzzz"

	assert.Panics(t, func() { IncrementKey([]byte(key), key[len(key)-1], len(key)-1) })
}

func TestIncrementKey_GivenKey_ShouldIncrement(t *testing.T) {
	key := "aaaaaaaaaaaaaaa"

	updatedKey := IncrementKey([]byte(key), key[len(key)-1], len(key)-1)
	assert.Equal(t, "aaaaaaaaaaaaaab", updatedKey)
}

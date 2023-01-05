package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
	TearDown()
}

func TearDown() {
	os.Remove(testConfigFileName)
}

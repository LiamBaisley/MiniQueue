package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	TestConfig.SecurityHash, hashErr = CreateHash("TestSecret")
	m.Run()
	TearDown()
}

func TearDown() {
	os.Remove(testConfigFileName)
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testConfigFileName = "testConfig.json"
var TestConfig Config
var hashErr error

func TestFileExist_GivenExistingFile_ShouldReturnTrue(t *testing.T) {
	filename := "main.go"

	exists := CheckFileExist(filename)

	assert.True(t, exists)
}

func TestFileExist_GivenNonExistingFile_ShouldReturnFalse(t *testing.T) {
	filename := "random.go"

	exists := CheckFileExist(filename)

	assert.False(t, exists)
}

func TestWriteConfig_ShouldWriteConfig(t *testing.T) {
	val := WriteConfig(TestConfig, testConfigFileName)

	assert.True(t, val)
	assert.FileExists(t, testConfigFileName)
}

func TestReadConfig_ShouldReadConfig(t *testing.T) {
	returnedConfig := GetConfig(testConfigFileName)

	assert.Equal(t, TestConfig, returnedConfig)
}

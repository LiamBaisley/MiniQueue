package main

import (
	"encoding/json"
	"errors"
	"os"
)

const ConfigFileName = "config.json"

type Config struct {
	SecurityHash string
}

func GetConfig() Config {
	content := ReadFile()

	config := Config{}
	err := json.Unmarshal(content, &config)

	if err != nil {
		panic("could not unmarshal json")
	}

	return config
}

func WriteConfig(config Config) bool {
	content, err := json.Marshal(config)
	if err != nil {
		panic("Could not marshal json")
	}

	result, err := WriteFile(content)
	if err != nil {
		panic("Could not write file")
	}

	return result
}

func ReadFile() []byte {
	content, err := os.ReadFile(ConfigFileName)
	if err != nil {
		panic("file does not exist")
	}

	return content
}

func WriteFile(content []byte) (bool, error) {
	err := os.WriteFile(ConfigFileName, content, 0644)
	if err != nil {
		return false, err
	}
	return true, nil
}

func CheckFileExist(filename string) bool {
	if _, err := os.Stat("./" + filename); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		panic("Could not determine if the file does or does not exist.")
	}
}

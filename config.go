package main

import (
	"encoding/json"
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

// Todo: switch to OS instead of ioutil
func ReadFile() []byte {
	content, err := os.ReadFile(ConfigFileName)
	if err != nil {
		panic("file does not exist")
	}

	return content
}

// Todo: switch to OS instead of ioutil
func WriteFile(content []byte) (bool, error) {
	err := os.WriteFile(ConfigFileName, content, 0644)
	if err != nil {
		return false, err
	}
	return true, nil
}

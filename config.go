package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Secret string
}

func GetEnvSecret() (string, error) {
	secret := os.Getenv("MINIQ-AUTH")
	if secret == "" {
		return "", fmt.Errorf("GetEnvSecret: no environment variable could be found.")
	}

	return secret, nil
}

func GetConfig(filename string) Config {
	content := ReadFile(filename)

	config := Config{}
	err := json.Unmarshal(content, &config)

	if err != nil {
		panic("Could not read configuration data.")
	}

	return config
}

func WriteConfig(config Config, filename string) bool {
	content, err := json.Marshal(config)
	if err != nil {
		panic("could not write configuration data")
	}

	result, err := WriteFile(content, filename)
	if err != nil || !result {
		panic("could not write configuration file")
	}
	fmt.Print(result)
	return result
}

func ReadFile(filename string) []byte {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic("file does not exist")
	}

	return content
}

func WriteFile(content []byte, filename string) (bool, error) {
	err := os.WriteFile(filename, content, 0644)
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

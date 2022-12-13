package main

import (
	"golang.org/x/crypto/bcrypt"
)

func CreateHash(value string) (string, error) {
	pass := []byte(value)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareHash(hash string, password string) (bool, error) {
	passByte := []byte(password)
	hashByte := []byte(hash)
	err := bcrypt.CompareHashAndPassword(hashByte, passByte)

	if err != nil {
		return false, err
	}

	return true, nil
}

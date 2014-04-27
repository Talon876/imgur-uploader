package main

import (
	"log"
	"math/rand"
)

const (
	ID_LENGTH = 8
)

var UrlCharacters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
var Symbols = []byte("`~!@#$%^&*()-=_+[]{}?")

//NewId generates a string that is the specified length and contains [A-Za-z0-9]
func NewId(length int) string {
	return generateString(length, UrlCharacters)
}

//NewIdKey generates a longer string containing more symbols.
func NewIdPassword() string {
	return generateString(ID_LENGTH*2, append(UrlCharacters, Symbols...))
}

//DefaultId generates a random id and password combination
func NewIdPair() (string, string) {
	return NewId(ID_LENGTH), NewIdPassword()
}

func generateString(length int, characters []byte) string {
	generatedId := make([]byte, length)
	for i := 0; i < length; i++ {
		generatedId[i] = characters[rand.Intn(len(characters))]
	}
	log.Println("Generated: " + string(generatedId))
	return string(generatedId)
}

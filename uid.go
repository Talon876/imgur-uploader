package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/go-martini/martini"
)

const (
	ID_LENGTH = 8
)

var UrlCharacters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
var Symbols = []byte("`~!@#$%^&*()-=_+[]{}?")

//NewId generates a string that is the specified length and contains [A-Za-z0-9]
func newId(length int) string {
	return generateString(length, UrlCharacters)
}

//newIdKey generates a longer string containing more symbols.
func newIdPassword() string {
	return generateString(ID_LENGTH*2, append(UrlCharacters, Symbols...))
}

//newIdPair generates a random id and password combination
func newIdPair() (string, string) {
	return newId(ID_LENGTH), newIdPassword()
}

func generateString(length int, characters []byte) string {
	generatedId := make([]byte, length)
	for i := 0; i < length; i++ {
		generatedId[i] = characters[rand.Intn(len(characters))]
	}
	log.Println("Generated: " + string(generatedId))
	return string(generatedId)
}

func ReserveRandomId() (int, string) {
	id, pass := newIdPair()
	return 200, fmt.Sprintf("%s:%s", id, pass)
}

func ReserveNamedId(params martini.Params) (int, string) {
	id, pass := params["id"], newIdPassword()
	return 200, fmt.Sprintf("%s:%s", id, pass)
}

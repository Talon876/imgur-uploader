package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Printf("Starting imgur mapping server...\n")
	m := martini.Classic()
	m.Get("/", displayIndex)
	m.Get("/hello/:name", greetings)
	m.Get("/:id", serveImage)
	m.Get("/map/:from/:to", addMapping)
	m.Post("/:id", receiveImage)
	m.Run()
}

func receiveImage(w http.ResponseWriter, req *http.Request, params martini.Params) {
	file, _, err := req.FormFile("file")
	if err != nil {
		fmt.Println(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile(params["id"]+".jpg", data, 0664)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Wrote " + params["id"] + ".jpg")
	file.Close()
}

func addMapping(res http.ResponseWriter, req *http.Request, params martini.Params) {
	urlMap[params["from"]] = params["to"]
	http.Redirect(res, req, "/"+params["from"], http.StatusTemporaryRedirect)
}

func serveImage(res http.ResponseWriter, req *http.Request, params martini.Params) {
	imageId := params["id"]
	fmt.Printf("Mapped %s -> %s\n", imageId, urlMap[imageId])
	http.Redirect(res, req, "http://i.imgur.com/"+urlMap[imageId], http.StatusMovedPermanently)
}

func displayIndex(params martini.Params) string {
	return generateImageId(ID_LENGTH)
}

func greetings(params martini.Params) string {
	return "Hello " + params["name"]
}

func generateImageId(length int) string {
	generatedId := make([]byte, length)
	for i := 0; i < length; i++ {
		generatedId[i] = UrlCharacters[rand.Intn(len(UrlCharacters))]
	}
	return string(generatedId)
}

var UrlCharacters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

var urlMap = map[string]string{}

const (
	ID_LENGTH = 8
)

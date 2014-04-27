package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-martini/martini"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Printf("Starting imgur mapping server...\n")
	m := martini.Classic()

	m.Get("/", displayHelp)

	m.Get("/id", reserveRandomId)
	m.Get("/id/:id", reserveId)

	m.Get("/img/:id", serveImage)
	m.Post("/img/:id", receiveImage)

	m.Run()
}

func displayHelp(params martini.Params) (int, string) {
	data, err := ioutil.ReadFile("help.txt")
	if err != nil {
		fmt.Println(err)
		return 404, "error getting help file"
	}
	return 200, string(data)
}

func reserveRandomId() string {
	generatedId, generatedPassword := generateImageId(ID_LENGTH), generateImageId(ID_LENGTH*2)
	fmt.Printf("Generated id:key combo; %s:%s\n", generatedId, generatedPassword)
	return fmt.Sprintf("%s:%s", generatedId, generatedPassword)
}

func reserveId(params martini.Params) (int, string) {
	reservedId, generatedPassword := params["id"], generateImageId(ID_LENGTH*2)
	fmt.Printf("Generated id:key combo; %s:%s\n", reservedId, generatedPassword)
	return 200, fmt.Sprintf("%s:%s", reservedId, generatedPassword)
}

func serveImage(res http.ResponseWriter, req *http.Request, params martini.Params) (int, []byte) {
	imageId := params["id"]
	data, err := ioutil.ReadFile(imageId + ".jpg")
	if err != nil {
		fmt.Println(err)
		return 404, []byte("error loading the " + imageId + " image")
	}
	return 200, data
}

func receiveImage(res http.ResponseWriter, req *http.Request, params martini.Params) {

	if imageKey := req.Header.Get("X-Mapper-Key"); len(imageKey) != 0 {
		fmt.Println("Found key '" + imageKey + "'")
	} else {
		fmt.Println("Could not find required key")
		res.Write([]byte("Could not find required X-Mapper-Key header"))
		return
	}

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
	file.Close()

	fmt.Println("Saved " + params["id"] + ".jpg")
	http.Redirect(res, req, "/img/"+params["id"], http.StatusTemporaryRedirect)
}

func generateImageId(length int) string {
	generatedId := make([]byte, length)
	for i := 0; i < length; i++ {
		generatedId[i] = UrlCharacters[rand.Intn(len(UrlCharacters))]
	}
	return string(generatedId)
}

var UrlCharacters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

const (
	ID_LENGTH = 8
)

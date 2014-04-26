package main

import (
	"crypto/rand"
	"fmt"
	"github.com/go-martini/martini"
	"io"
	"net/http"
	"os"
)

func main() {
	fmt.Printf("Starting imgur server...\n")
	populateMap()
	m := martini.Classic()
	m.Get("/", displayIndex)
	m.Get("/hello/:name", greetings)
	m.Get("/:id", serveImage)
	m.Get("/:from/:to", addMapping)
	m.Post("/:id", receiveImage)
	m.Run()
}

func receiveImage(res http.ResponseWriter, req *http.Request, params martini.Params) {
	fmt.Printf("Uploading %s\n", params["id"])
	file, _, err := req.FormFile("file")

	fmt.Fprintf(res, "/%s.jpg\n", params["id"])

	defer file.Close()

	if err != nil {
		fmt.Fprintln(res, err)
		return
	}

	out, err := os.Create(params["id"] + ".jpg")
	if err != nil {
		fmt.Fprintf(res, "Failed to open the file for writing")
		return
	}

	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Print(res, "Couldn't copy: ")
		fmt.Fprintln(res, err)
	}

	fmt.Fprintf(res, "File %s.jpg uploaded successfully.", params["id"])
}

func addMapping(res http.ResponseWriter, req *http.Request, params martini.Params) {
	urlMap[params["from"]] = params["to"]
	http.Redirect(res, req, "/"+params["from"], http.StatusTemporaryRedirect)
}

func populateMap() {
	urlMap["test"] = "k1fctS9.jpg"
	urlMap["abc"] = "vHV1FRg.jpg"
	urlMap["def"] = "5na6c1v.jpg"
	urlMap["ghi"] = "lce0i.jpg"
	urlMap["jkl"] = "hJ6Sx.jpg"
	urlMap["mno"] = "w3naz.jpg"
	urlMap["pqr"] = "kkyv9.jpg"
	urlMap["stu"] = "9mlnw.jpg"
	urlMap["vwx"] = "dSpzF.gif"
	urlMap["yz"] = "kkyv9.jpg"
}

func serveImage(res http.ResponseWriter, req *http.Request, params martini.Params) {
	imageId := params["id"]
	fmt.Printf("Mapped %s -> %s\n", imageId, urlMap[imageId])
	http.Redirect(res, req, "http://i.imgur.com/"+urlMap[imageId], http.StatusMovedPermanently)
}

func displayIndex(params martini.Params) string {
	return generateImageId()
}

func greetings(params martini.Params) string {
	return "Hello " + params["name"]
}

func generateImageId() string {
	bytes := make([]byte, URL_LENGTH)
	randomBytes := make([]byte, 10)
	maxrb := byte(256 - (256 % len(UrlCharacters)))
	i := 0

	for {
		if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
			panic("Error reading from random source: " + err.Error())
		}

		for _, c := range randomBytes {
			if c >= maxrb {
				continue
			}

			bytes[i] = UrlCharacters[c%byte(len(UrlCharacters))]
			i++
			if i == len(bytes) {
				return string(bytes)
			}
		}
	}
	return string(randomBytes)
}

var UrlCharacters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

const (
	URL_LENGTH = 8
)

var urlMap = map[string]string{}

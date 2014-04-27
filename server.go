package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/go-martini/martini"
)

var imgurKey string

func init() {
	flag.StringVar(&imgurKey, "imgurkey", "ENTERYOURKEY", "Your imgur v3 api client id")
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()
	fmt.Println("Starting imgur mapping server...")
	fmt.Println("Using api key: " + imgurKey)
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
		http.Error(res, "Could not find required X-Mapper-Key header", http.StatusBadRequest)
		return
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		fmt.Println("Could not find 'file' parameter.")
		http.Error(res, "Could not find 'file' parameter", http.StatusBadRequest)
		return
	}

	fmt.Println("Size: " + req.Header.Get(""))

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "internal error", http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile(params["id"]+".jpg", data, 0664)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "internal error", http.StatusInternalServerError)
		return
	}
	file.Close()

	fmt.Println("Saved " + params["id"] + ".jpg")
	imgurResponse, err := sendToImgur(params["id"])
	if err != nil {
		fmt.Println("Failed to upload to imgur")
		http.Redirect(res, req, "/img/"+params["id"], http.StatusTemporaryRedirect)
	}
	fmt.Println(params["id"] + " -> " + imgurResponse.Data.Link)
	//http.Redirect(res, req, imgurResponse.Data.Link, http.StatusTemporaryRedirect)
	res.Write([]byte(imgurResponse.Data.Link))

}

func sendToImgur(imageId string) (*ImgurResponse, error) {

	var imageBytes bytes.Buffer
	imageWriter := multipart.NewWriter(&imageBytes)

	imageFile, err := os.Open(imageId + ".jpg")

	if err != nil {
		fmt.Println("Couldn't open file " + imageId + ".jpg")
		return nil, errors.New("Couldn't open file " + imageId + ".jpg")
	}

	formWriter, err := imageWriter.CreateFormFile("image", imageId+".jpg")
	if err != nil {
		fmt.Println("Couldn't create file form")
		return nil, errors.New("Couldn't create file form")
	}

	if _, err = io.Copy(formWriter, imageFile); err != nil {
		fmt.Println("Couldn't copy image data in to form")
		return nil, errors.New("Couldn't copy image data in to form")
	}

	imageWriter.Close()

	req, err := http.NewRequest("POST", IMGUR_API_ENDPOINT, &imageBytes)

	if err != nil {
		fmt.Println("Couldn't create http request")
		return nil, errors.New("Couldn't create http request")
	}

	req.Header.Set("Content-Type", imageWriter.FormDataContentType())
	req.Header.Set("Authorization", "Client-ID "+imgurKey)

	fmt.Println("Uploading " + imageId + ".jpg to imgur...")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Did not get a response from imgur")
		return nil, errors.New("Did not get a response from imgur")
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Could not read body from imgur")
		return nil, errors.New("Could not read body from imgur")
	}

	var imgurResponse ImgurResponse
	err = json.Unmarshal(body, &imgurResponse)
	if err != nil {
		fmt.Println("Couldn't decode json response from imgur: " + string(body))
	}
	fmt.Println("Done uploading " + imageId)
	return &imgurResponse, nil
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
	ID_LENGTH            = 8
	IMGUR_API_ENDPOINT   = "https://api.imgur.com/3/image"
	MAX_IMAGE_SIZE_BYTES = 1024 * 5
)

type ImgurResponse struct {
	Data *ImgurData
}

type ImgurData struct {
	Id         string
	Deletehash string
	Link       string
}

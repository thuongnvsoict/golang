package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Image struct{
	Publisher 	string 	`json:"Publisher"`
	Url			string 	`json:"Url"`
	Status 		int 	`json:"status"`
	IsThumb 	int 	`json:"isThumb"`
	Error 		string 	`json:"error"`
}

type ImageResultDownload struct{
	Url			string
	FileName	string
	Error		error
	Status		string
}

const INPUT_IMAGE_DIR = "./image_input.json"
const OUTPUT_IMAGE_DIR = "/home/thuongnv/ImageDownloads/"

/*----------------------------------------------
		DOWNLOAD IMAGES USING HTTP GET
/---------------------------------------------*/
func downloadImage(url string, OutputDestination string, imageChan chan ImageResultDownload) {

	// GET FILE NAME
	filename := filepath.Base(url)

	// 	DOWNLOAD IMAGE
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
	}
	defer res.Body.Close()

	out, err := os.Create(OutputDestination + filename)
	if err != nil {
		fmt.Printf("%s", err)
	}
	defer out.Close()
	_, err = io.Copy(out, res.Body)

	//	PUSH ImageResultDownload TO CHANNELS
	if err != nil{
		fmt.Printf("%s", err)
		imageChan <- ImageResultDownload{url, filename, nil , "Failed" }
	} else{
		imageChan <- ImageResultDownload{url, filename, err, "Successul"}
	}
}

func importImageFromJson(directory string) []Image{

	var images []Image
	/*--------------------------
		READ IMAGES JSON FILE
	 -------------------------*/
	// Open our jsonFile
	jsonImageFile, err := os.Open(directory)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened " + directory)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonImageFile.Close()
	byteImageValue, _ := ioutil.ReadAll(jsonImageFile)

	json.Unmarshal([]byte(byteImageValue), &images)

	return images
}


func downloadImages(images []Image, outputDestination string){

	/*-------------------------------------------
		DOWNLOAD VIDEOS AND THEIR THUMBNAIL
	 ------------------------------------------*/
	//Declare channels
	imageChan := make(chan ImageResultDownload)

	// Download Images
	numberOfImage := 0
	for i, image := range images{
		if image.Error == ""{
			numberOfImage++
			fmt.Printf("%d: %s \n", i, image.Url)
			go downloadImage(image.Url, outputDestination, imageChan)
		}
	}

	//Getting data from channels
	for j := 0; j < numberOfImage; j++{
		msgImage := <-imageChan
		fmt.Printf("%100s :%s - %s\n", msgImage.FileName, msgImage.Status)
	}
}


func main(){
	var images []Image
	images = importImageFromJson(INPUT_IMAGE_DIR)

	startImage := time.Now()

	downloadImages(images, OUTPUT_IMAGE_DIR)

	elapsedImage := time.Since(startImage)
	fmt.Printf("Download images took %s\n", elapsedImage)

}
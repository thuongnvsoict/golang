package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	INPUT_IMAGE_DIR 	= "./image_input.json"
	InputVideoDir    	= "./input.json"
)

func importVideoFromJson(directory string) []Video {

	var videos []Video
	/*--------------------------
		READ VIDEOS JSON FILE
	-------------------------*/
	// Open our jsonFile
	jsonVideoFile, err := os.Open(directory)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened " + directory)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonVideoFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonVideoFile)

	json.Unmarshal([]byte(byteValue), &videos)

	return videos
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

	err = json.Unmarshal([]byte(byteImageValue), &images)
	if err != nil {
		return nil
	}

	return images
}

func printResult(downloadedVideo []DownloadedVideo, downloadedImage []DownloadedImage){
	for _,videoResult := range downloadedVideo{
		fmt.Printf("%100s :%d - %d\n", videoResult.VideoFileName, videoResult.VideoDownloadedStatus, videoResult.ThumbDownloadedStatus)
	}
	for _,imageResult := range downloadedImage{
		fmt.Printf("%100s :%d\n", imageResult.Url, imageResult.Status)
	}
}

func main(){

	//1. retrieve data from database, api or RabbitMQ
	// Download Video
	var videos []Video
	var images []Image

	videos = importVideoFromJson(InputVideoDir)
	images = importImageFromJson(INPUT_IMAGE_DIR)

	//fmt.Printf("%d %d\n", len(videos), len(images))

	// 2. push download items to a collection to download
	downloadManager := NewDownloadManager(OutputDownloadVideoDir, OutputDownloadImageDir)

	downloadManager.PushImages(images[0:6])
	downloadManager.PushVideos(videos[0:6])

	startVideo := time.Now()

	// downloadedVideo, downloadedImage := downloadManager.DownloadAll(OUTPUT_VIDEO_DIR, OUTPUT_IMAGE_DIR)
	downloadedVideo, downloadedImage := downloadManager.DownloadAll()

	elapsedVideo := time.Since(startVideo)
	fmt.Printf("Download videos took %s\n", elapsedVideo)

	printResult(downloadedVideo, downloadedImage)

	//3. upload object to the server
}
package main

import (
	"fmt"
	"time"
)


func main(){

	// Download Video
	var videos []Video
	videos = importVideoFromJson(INPUT_VIDEO_DIR)

	startVideo := time.Now()

	downloadVideos(videos, OUTPUT_VIDEO_DIR)

	elapsedVideo := time.Since(startVideo)
	fmt.Printf("Download videos took %s\n", elapsedVideo)

	// Download Image
	var images []Image
	images = importImageFromJson(INPUT_IMAGE_DIR)

	startImage := time.Now()

	downloadImages(images, OUTPUT_IMAGE_DIR)

	elapsedImage := time.Since(startImage)
	fmt.Printf("Download images took %s\n", elapsedImage)
}

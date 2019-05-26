package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type DownloadImageManager struct {
	outputDirectory string
}

/**
 * Returns a new DownloadVideoManager object
 */
func NewDownloadImageManager(outputDir string) DownloadImageManager {
	return DownloadImageManager {
		outputDirectory: outputDir,
	}
}

/**
 *		DOWNLOAD IMAGES USING HTTP GET
*/
func (imageManager *DownloadImageManager) downloadImage(url string, imageChan chan DownloadedImage) {

	// GET FILE NAME
	fileName := filepath.Base(url)

	// 	DOWNLOAD IMAGE
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
	}
	defer res.Body.Close()

	out, err := os.Create(imageManager.outputDirectory + fileName)
	if err != nil {
		fmt.Printf("%s", err)
	}
	defer out.Close()
	_, err = io.Copy(out, res.Body)

	//	PUSH ImageResultDownload TO CHANNELS
	if err != nil{
		fmt.Printf("%s", err)
		imageChan <- DownloadedImage{url, StatusDownloadImageFailed , err }
	} else{
		imageChan <- DownloadedImage{url, StatusDownloadSuccess, err}
	}
}
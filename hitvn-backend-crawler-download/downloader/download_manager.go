package main

import (
	"strings"
)

const (
	VideoTypeMp4 = ".mp4"
	VideoTypeM3U8 = ".m3u8"
	VideoTypeYouTube = "youtube"
	VideoTypeVimeo = "vimeo"
)

type DownloadManager struct {
	downloadCollection   DownloadCollection
	downloadVideoManager DownloadVideoManager
	downloadImageManager DownloadImageManager
}

/**
 * Returns a new DownloadManager
 */
func NewDownloadManager(outputVideoDir string, outputImageDir string) DownloadManager {
	dVideoManager := NewDownloadVideoManager(outputVideoDir)
	dImageManager := NewDownloadImageManager(outputImageDir)

	return DownloadManager{
		downloadVideoManager	: dVideoManager,
		downloadImageManager 	: dImageManager,
	}
}

/**
 * Pushes image to download
 */
func (manager *DownloadManager) PushImage(image Image) {
	manager.downloadCollection.images = append(manager.downloadCollection.images, image)

}

/**
 * Pushes images to download
 */
func (manager *DownloadManager) PushImages(images []Image) {
	manager.downloadCollection.images = append(manager.downloadCollection.images, images...)
}

/**
 * Pushes video to download
 */
func (manager *DownloadManager) PushVideo(video Video) {
	manager.downloadCollection.videos = append(manager.downloadCollection.videos, video)
}

/**
 * Pushes videos to download
 */
func (manager *DownloadManager) PushVideos(videos []Video) {
	manager.downloadCollection.videos = append(manager.downloadCollection.videos, videos...)
}

/**
 * Downloads all media
 */
func (manager *DownloadManager) DownloadAll() ([]DownloadedVideo, []DownloadedImage) {
	// Declare videos channels
	mp4Chan 	:= make(chan DownloadedVideo)
	m3u8Chan 	:= make(chan DownloadedVideo)
	youtubeChan := make(chan DownloadedVideo)

	// Declare images  channels
	imageChan := make(chan DownloadedImage)

	// Download videos
	numberOfVideo := 0
	for _, video := range manager.downloadCollection.videos {
		if len(video.Error) <= 0 {
			numberOfVideo++
			var flag string
			if strings.Contains(video.Url, ".mp4"){
				flag = VideoTypeMp4
			} else if strings.Contains(video.Url, ".m3u8") {
				flag = VideoTypeM3U8
			} else if strings.Contains(video.Url, "youtube.com/") {
				flag = VideoTypeYouTube
			} else if strings.Contains(video.Url, "vimeo.com/") {
				flag = VideoTypeVimeo
			} else {
				flag = "unknown"
			}

			fileName := GetFileName(video.Url, flag)

			switch flag {
			case VideoTypeMp4:
				go manager.downloadVideoManager.DownloadMP4(video.Url, fileName, mp4Chan)
			case VideoTypeM3U8:
				go manager.downloadVideoManager.DownloadFFMPEG(video.Url, video.PostUrl, fileName, m3u8Chan)
			case VideoTypeYouTube:
				go manager.downloadVideoManager.DownloadYoutube(video.Url, youtubeChan)
			case VideoTypeVimeo:
				go manager.downloadVideoManager.DownloadYoutube(video.Url, youtubeChan)
			}
		}
	}

	// Download images
	numberOfImage := 0
	for _, image := range manager.downloadCollection.images{
		if image.Error == ""{
			numberOfImage++
			//fmt.Printf("%d: %s \n", i, image.Url)
			go manager.downloadImageManager.downloadImage(image.Url, imageChan)
		}
	}

	// Getting data from channels
	var downloadedImages []DownloadedImage
	var downloadedVideos []DownloadedVideo

	// Set total download files
	total := numberOfVideo + numberOfImage

	for i := 0; i < total; i++{
		select{
		case msgVideo1 := <-mp4Chan:
			downloadedVideos = append(downloadedVideos, msgVideo1)
			//fmt.Printf("%100s :%d - %d\n", msgVideo1.VideoFileName, msgVideo1.VideoDownloadedStatus, msgVideo1.ThumbDownloadedStatus)
		case msgVideo2 := <-m3u8Chan:
			downloadedVideos = append(downloadedVideos, msgVideo2)
			//fmt.Printf("%100s :%d - %d\n", msgVideo2.VideoFileName, msgVideo2.VideoDownloadedStatus, msgVideo2.ThumbDownloadedStatus)
		case msgVideo3 := <-youtubeChan:
			downloadedVideos = append(downloadedVideos, msgVideo3)
			//fmt.Printf("%100s :%d - %d\n", msgVideo3.VideoFileName, msgVideo3.VideoDownloadedStatus, msgVideo3.ThumbDownloadedStatus)
		case msgImage := <-imageChan:
			downloadedImages = append(downloadedImages, msgImage)
			//fmt.Printf("%100s :%d\n", msgImage.Url, msgImage.Status)
		}
	}

	return downloadedVideos, downloadedImages
}
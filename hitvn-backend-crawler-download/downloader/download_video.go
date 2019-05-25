package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// STRUCT INPUT
type Video struct{
	Publisher 	string 	`json:"Publisher"`
	Url			string 	`json:"Url"`
	Status 		int 	`json:"status"`
	Error 		string 	`json:"error"`
	Posturl 	string 	`json:"posturl"`
}

// STRUCT RESULT DOWNLOAD
type VideoResultDownload struct{
	Url			string
	FileName	string
	Error		error
	Status		string
	ThumbError	error
	ThumbStatus	string
}

const (
	INPUT_VIDEO_DIR = "../input.json"
	OUTPUT_VIDEO_DIR = "/home/thuongnv/VideoDownloads/"
	VIDEO_THUMB_TIME = "00:00:03"
	VIDEO_THUMB_SIZE = "480x320"
)

/*--------------------------------------
		GET FILE NAME FROM URL
/-------------------------------------*/
func getFileName(url string, videoType string) string{

	filename := url[:strings.Index(url, videoType)]
	filename = filepath.Base(filename)
	filename = filename + ".mp4"
	return filename

}

/*----------------------------------------------
		DOWNLOAD MP4 VIDEOS USING HTTP GET
/---------------------------------------------*/
func downloadMP4(url string, filename string, outputDirectory string, mp4Chan chan VideoResultDownload) {

	// 	DOWNLOAD VIDEO
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
	}
	defer res.Body.Close()

	out, err := os.Create(outputDirectory + filename)
	if err != nil {
		fmt.Printf("%s", err)
	}
	defer out.Close()
	_, err = io.Copy(out, res.Body)

	//	DOWNLOAD THUMBNAIL
	thumbFileName := filename[:strings.Index(filename, ".mp4")] + ".jpg"
	_, errThumb := exec.Command("ffmpeg", "-i", outputDirectory + filename, "-ss", VIDEO_THUMB_TIME ,"-vframes", "1", "-s", VIDEO_THUMB_SIZE, outputDirectory + thumbFileName).CombinedOutput()

	//	PUSH VideoResultDownload TO CHANNELS
	if err != nil{
		fmt.Printf("%s", err)
		mp4Chan <- VideoResultDownload{url, filename, nil , "Failed", errThumb, "None"}
	} else{
		if errThumb != nil{
			mp4Chan <- VideoResultDownload{url, filename, err, "Successul", errThumb,"Failed"}
		}else{
			mp4Chan <- VideoResultDownload{url, filename, err, "Successul", errThumb,"Successul"}
		}
	}
}

/*-------------------------------------------
	DOWNLOAD STREAM VIDEOS USING FFMPEG
-------------------------------------------*/
func downloadFFMPEG(url string, postUrl string, filename string, outputDirectory string,  m3u8Chan chan VideoResultDownload){

	// 	DOWNLOAD VIDEO
	//	Syntax: ffmpeg -headers "Referer: <URL_Post>" -i “<url>” -c copy -bsf:a aac_adtstoasc <output_file_name>
	//  -y: overide file
	_, err := exec.Command("ffmpeg", "-y","-headers", "Referer: "+ postUrl , "-i" ,url, "-c", "copy", "-bsf:a", "aac_adtstoasc",outputDirectory + filename).CombinedOutput()

	//	DOWNLOAD THUMBNAIL
	// 	Syntax: ffmpeg -i "$f" -ss 00:00:03 -vframes 1 -s 480x320 "${f%.mp4}.jpg"
	thumbFileName := filename[:strings.Index(filename, ".mp4")] + ".jpg"
	_, errThumb := exec.Command("ffmpeg", "-i", outputDirectory + filename, "-ss", VIDEO_THUMB_TIME ,"-vframes", "1", "-s", VIDEO_THUMB_SIZE, outputDirectory + thumbFileName).CombinedOutput()

	//	PUSH VideoResultDownload TO CHANNELS
	if err != nil{
		fmt.Printf("%s", err)
		m3u8Chan <- VideoResultDownload{url, filename, nil , "Failed", errThumb,"None"}
	} else{
		if errThumb != nil{
			m3u8Chan <- VideoResultDownload{url, filename, err, "Successul", errThumb,"Failed"}
		}else{
			m3u8Chan <- VideoResultDownload{url, filename, err, "Successul", errThumb,"Successul"}
		}
	}
}

/*----------------------------------------------
	DOWNLOAD YOUTUBE VIDEOS USING YOUTUBE-DL
----------------------------------------------*/
func downloadYoutube(url string, outputDirectory string, youtubeChan chan VideoResultDownload){

	//	DOWNLOAD VIDEOS AND THEIR THUMBNAIL
	_, err := exec.Command("youtube-dl", "-o", outputDirectory + "%(title)s", url, "--write-thumbnail").CombinedOutput()

	//	PUSH VideoResultDownload TO CHANNELS
	if err != nil{
		fmt.Printf("%s", err)
		youtubeChan <- VideoResultDownload{url, url, nil , "Failed", err, "None"}
	}else{
		youtubeChan <- VideoResultDownload{url, url, err, "Successul", err, "Successul"}
	}
}

func importVideoFromJson(directory string) []Video{

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


func downloadVideos(videos []Video, outputDirectory string){
	/*-------------------------------------------
		DOWNLOAD VIDEOS AND THEIR THUMBNAIL
	------------------------------------------*/
	//Declare channels
	mp4Chan := make(chan VideoResultDownload)
	m3u8Chan := make(chan VideoResultDownload)
	youtubeChan := make(chan VideoResultDownload)

	// Download videos
	numberOfVideo := 0
	for i, video := range videos[:6]{
		if video.Error == ""{
			numberOfVideo++
			var flag string
			if strings.Contains(video.Url, ".mp4"){
				flag = ".mp4"
			}else if strings.Contains(video.Url, ".m3u8"){
				flag = ".m3u8"
			}else if strings.Contains(video.Url, "youtube.com/"){
				flag = "youtube"
			}else if strings.Contains(video.Url, "vimeo.com/"){
				flag = "vimeo"
			}else{
				flag = "unknow"
			}

			fmt.Printf("%d: %s - %s\n", i, flag, video.Url)

			filename := getFileName(video.Url, flag)

			switch flag {
			case ".mp4":
				go downloadMP4(video.Url, filename, outputDirectory,mp4Chan)
			case ".m3u8":
				go downloadFFMPEG(video.Url, video.Posturl, filename, outputDirectory, m3u8Chan)
			case "youtube":
				go downloadYoutube(video.Url, outputDirectory, youtubeChan)
			case "vimeo":
				go downloadYoutube(video.Url, outputDirectory, youtubeChan)
			}
		}

	}

	//Getting data from channels
	for i := 0; i < numberOfVideo; i++{
		select{
		case msg1 := <-mp4Chan:
			fmt.Printf("%100s :%s - %s\n", msg1.FileName, msg1.Status, msg1.ThumbError)
		case msg2 := <-m3u8Chan:
			fmt.Printf("%100s :%s - %s\n", msg2.FileName, msg2.Status, msg2.ThumbError)
		case msg3 := <-youtubeChan:
			fmt.Printf("%100s :%s - %s\n", msg3.FileName, msg3.Status, msg3.ThumbError)
		}
	}

}


func main(){
	var videos []Video
	videos = importVideoFromJson(INPUT_VIDEO_DIR)

	startVideo := time.Now()

	downloadVideos(videos, OUTPUT_VIDEO_DIR)

	elapsedVideo := time.Since(startVideo)
	fmt.Printf("Download videos took %s\n", elapsedVideo)

}


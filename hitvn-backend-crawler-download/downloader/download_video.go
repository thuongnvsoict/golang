package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type DownloadVideoManager struct {
	outputDirectory string
}

/**
 * Returns a new DownloadVideoManager object
 */
func NewDownloadVideoManager(outputDir string) DownloadVideoManager {
	return DownloadVideoManager {
		outputDirectory: outputDir,
	}
}

/**
 * Downloads a mp4 video
 */
func (videoManager *DownloadVideoManager) DownloadMP4(url string, filename string, mp4Chan chan DownloadedVideo) {
	// 	DOWNLOAD VIDEO
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
	}
	defer res.Body.Close()

	out, err := os.Create(videoManager.outputDirectory + filename)
	if err != nil {
		fmt.Printf("%s", err)
	}
	defer out.Close()
	_, err = io.Copy(out, res.Body)

	//	DOWNLOAD THUMBNAIL
	thumbFileName := filename[:strings.Index(filename, ".mp4")] + ".jpg"
	_, errThumb := exec.Command("ffmpeg", "-i", videoManager.outputDirectory + filename, "-ss", VideoThumbTime,"-vframes", "1", videoManager.outputDirectory + thumbFileName).CombinedOutput()

	//	PUSH DownloadedVideo TO CHANNELS
	if err != nil{
		fmt.Printf("%s", err)
		mp4Chan <- DownloadedVideo {
			url,
			filename,
			StatusDownloadVideoFailed,
			err,
			thumbFileName,
			StatusDownloadVideoThumbFailed,
			err,
		}
	} else{
		if errThumb != nil{
			mp4Chan <- DownloadedVideo {
				url,
				filename,
				StatusDownloadSuccess,
				nil,
				thumbFileName,
				StatusDownloadVideoThumbFailed,
				err}
		}else{
			mp4Chan <- DownloadedVideo {
				url,
				filename,
				StatusDownloadSuccess ,
				nil,
				thumbFileName,
				StatusDownloadSuccess,
				nil,
			}
		}
	}
}

/**
 * Downloads a m3u8 video using FFmpeg
 */
func (videoManager *DownloadVideoManager) DownloadFFMPEG(url string, postUrl string, filename string, m3u8Chan chan DownloadedVideo){
	// 	DOWNLOAD VIDEO
	//	Syntax: ffmpeg -headers "Referer: <URL_Post>" -i “<url>” -c copy -bsf:a aac_adtstoasc <output_file_name>
	//  -y: overide file
	_, err := exec.Command("ffmpeg", "-y","-headers", "Referer: "+ postUrl , "-i" ,url, "-c", "copy", "-bsf:a", "aac_adtstoasc ", videoManager.outputDirectory + filename).CombinedOutput()

	//	DOWNLOAD THUMBNAIL
	// 	Syntax: ffmpeg -i "$f" -ss 00:00:03 -vframes 1 -s 480x320 "${f%.mp4}.jpg"
	thumbFileName := filename[:strings.Index(filename, ".mp4")] + ".jpg"
	_, errThumb := exec.Command("ffmpeg", "-i", videoManager.outputDirectory + filename, "-ss", VideoThumbTime,"-vframes", "1", videoManager.outputDirectory  + thumbFileName).CombinedOutput()

	//	PUSH DownloadedVideo TO CHANNELS
	if err != nil{
		fmt.Printf("%s", err)
		m3u8Chan <- DownloadedVideo {
			url,
			filename,
			StatusDownloadVideoFailed,
			err,
			thumbFileName,
			StatusDownloadVideoThumbFailed,
			err,
		}
	} else{
		if errThumb != nil{
			m3u8Chan <- DownloadedVideo{
				url,
				filename,
				StatusDownloadSuccess,
				nil,
				thumbFileName,
				StatusDownloadVideoThumbFailed,
				err,
			}
		}else{
			m3u8Chan <- DownloadedVideo{
				url,
				filename,
				StatusDownloadSuccess ,
				nil,
				thumbFileName,
				StatusDownloadSuccess,
				nil,
			}
		}
	}
}

/**
 * Downloads a youtube video
 */
func (videoManager *DownloadVideoManager) DownloadYoutube(url string, youtubeChan chan DownloadedVideo) {
	//	DOWNLOAD VIDEOS AND THEIR THUMBNAIL
	_, err := exec.Command("youtube-dl", "-o", videoManager.outputDirectory  + "%(title)s", url, "--write-thumbnail").CombinedOutput()

	//	PUSH DownloadedVideo TO CHANNELS
	if err != nil{
		fmt.Printf("%s", err)
		youtubeChan <- DownloadedVideo{
			url,
			url,
			StatusDownloadVideoFailed ,
			err,
			url,
			StatusDownloadVideoFailed,
			err,
		}
	}else{
		youtubeChan <- DownloadedVideo{
			url,
			url,
			StatusDownloadSuccess,
			nil,
			url,
			StatusDownloadSuccess,
			nil,
		}
	}
}
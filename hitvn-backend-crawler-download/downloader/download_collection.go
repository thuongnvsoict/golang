package main


const (
	OutputDownloadVideoDir = "/home/thuongnv/VideoDownloads/"
	OutputDownloadImageDir = "/home/thuongnv/ImageDownloads/"
	VideoThumbTime    = "00:00:03"
)

const (
	StatusDownloadSuccess 			int = 800
	StatusDownloadImageFailed 		int = 801
	StatusDownloadVideoFailed 		int = 802
	StatusDownloadVideoThumbFailed 	int = 803
)

type DownloadCollection struct {
	images []Image
	videos []Video
}

type Image struct {
	Url     string
	Status  int
	IsThumb int
	Error   string
}

type DownloadedImage struct {
	Url    string
	Status int
	Error  error
}

type Video struct {
	Url     string
	PostUrl string
	Status  int
	Error   string
}

type DownloadedVideo struct {
	Url                   string
	VideoFileName         string
	VideoDownloadedStatus int
	VideoError            error
	ThumbFileName         string
	ThumbDownloadedStatus int
	ThumbError            error
}
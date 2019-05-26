package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)
const OUTPUT_DIRECTORY = "/home/thuongnv/VideoDownloads/"

func downloadMP4Test(url string, filepath string) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
	}

	defer res.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("%s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)


}

func downloadFFMPEGTest(url string, destination string){

	//	ffmpeg -headers "Referer: <URL_Post>" -i “<url>” -c copy -bsf:a aac_adtstoasc <output_file_name>
	//  -y: overide file
	//_, err := exec.Command("ffmpeg", "-y","-headers", "Referer: "+ postUrl , "-i" ,url, "-c", "copy", "-bsf:a", "aac_adtstoasc",destination).CombinedOutput()
	_, err := exec.Command("ffmpeg", "-i", url, destination).CombinedOutput()
	if err != nil {
		fmt.Printf("%s", err)
	}
}

func getFileNameTest(url string, videoType string) string{

	if videoType == ".m3u8" {
		if strings.Contains(url, "/d1.vnecdn.net/"){
			var re = regexp.MustCompile(`(?m)[a-z0-9\-\_]+\-\d{8,}`)
			match := re.FindString(url)

			return  match + ".mp4"
		}else if strings.Contains(url, "ss-hls.catscdn.vn") {
			var re = regexp.MustCompile(`(?m)[a-z0-9\-\_]+\.mp4`)
			match := re.FindString(url)

			return  match + ".mp4"
		}
	}

	filename := url[:strings.Index(url, videoType)]
	filename = filepath.Base(filename)
	filename = filename + ".mp4"

	return filename
}


func main()  {
	//		https://d1.vnecdn.net/giaitri/video/video/web/mp4/2017/04/12/truong-thi-may-mua-dieu-khmer-1491966555/index-v1-a1.m3u8
	url := "https://d1.vnecdn.net/ngoisao/video/video/web/mp4/,240p,360p,480p,,/2019/05/24/trailer-ca-si-than-tuong-1558687965/vne/master.m3u8"
	fmt.Printf("%s", getFileNameTest(url, ".m3u8"))

	//url := "https://hls.mediacdn.vn/afamily/2019/5/22/0046491293-1558539735217597574749-4a372.mp4"
	//start := time.Now()
	//downloadMP4Test(url, OUTPUT_DIRECTORY +  "xxx.mp4")
	//elapsed := time.Since(start)
	//fmt.Printf("Binomial took %s\n", elapsed)
	//
	//start2 := time.Now()
	//downloadFFMPEGTest(url, OUTPUT_DIRECTORY +  "yyy.mp4")
	//elapsed2 := time.Since(start2)
	//fmt.Printf("Binomial took %s\n", elapsed2)

	//filename := "master.mp4"
	//thumbFileName := filename[:strings.Index(filename, ".mp4")] + ".jpg"
	//fmt.Println(thumbFileName)
}

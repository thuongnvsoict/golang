package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func DownloadFile(filepath string, url string) error{
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	return err
}

func executeCommand(program string, input string, output string){
	stdout, err := exec.Command(program, "-i", "-c", "copy", "-bsf:a", "aac_adtstoasc" ,input ,output).CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(stdout))
}

func main()  {

	//executeCommand("ffmpeg", "https://ss-hls.catscdn.vn/2019/05/18/5217819/than-tuong-bolero-2019-tap-7-cong-bo-top-5-doi-hlv-ngoc-son-giang-hong-ngoc.mp4/index.m3u8", "./hello.mp4")

	//https://f9-stream.nixcdn.com/PreNCT16/DungYeuNuaEmMetRoi-MIN-5946266.mp4?st=K7Fuy6AE4M-dX80MHpZacg&e=1558432110&t=1558345711060
	//https://mnmedias.api.telequebec.tv/m3u8/29880.m3u8

	fileUrl := "https://d1.vnecdn.net/vnexpress/video/video/web/mp4/,240p,360p,480p,,/2019/05/19/bai-tap-the-hinh-mien-phi-duoi-chan-cau-long-bien-1558263377/vne/master.m3u8"

	if err := DownloadFile("vnexpress.m3u8", fileUrl); err != nil {
		panic(err)
	}


	//// Instantiate default collector
	//c := colly.NewCollector(
	//	// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
	//	colly.AllowedDomains("hackerspaces.org"),
	//)
	//
	//// On every a element which has href attribute call callback
	//c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	//	link := e.Attr("href")
	//	// Print link
	//	fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	//	// Visit link found on page
	//	// Only those links are visited which are in AllowedDomains
	//	c.Visit(e.Request.AbsoluteURL(link))
	//})
	//
	//// Before making a request print "Visiting ..."
	//c.OnRequest(func(r *colly.Request) {
	//	fmt.Println("Visiting", r.URL.String())
	//})
	//
	//// Start scraping on https://hackerspaces.org
	//c.Visit("https://hackerspaces.org/")


//	// create client
//	client := grab.NewClient()
//	req, _ := grab.NewRequest(".", "https://ss-hls.catscdn.vn/2019/05/17/5207487/sq-1.mp4/index.m3u8")
//
//	// start download
//	fmt.Printf("Downloading %v...\n", req.URL())
//	resp := client.Do(req)
//	fmt.Printf("  %v\n", resp.HTTPResponse.Status)
//
//	// start UI loop
//	t := time.NewTicker(500 * time.Millisecond)
//	defer t.Stop()
//
//Loop:
//	for {
//		select {
//		case <-t.C:
//			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
//				resp.BytesComplete(),
//				resp.Size,
//				100*resp.Progress())
//
//		case <-resp.Done:
//			// download is complete
//			break Loop
//		}
//	}
//
//	// check for errors
//	if err := resp.Err(); err != nil {
//		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
//		os.Exit(1)
//	}
//
//	fmt.Printf("Download saved to ./%v \n", resp.Filename)


}

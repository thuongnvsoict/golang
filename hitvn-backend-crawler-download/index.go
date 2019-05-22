package main

import (
	"flag"
	"fmt"
	. "github.com/kkdai/youtube"
	"log"
	"os/user"
	"path/filepath"
)

func main(){
	flag.Parse()
	log.Println(flag.Args())
	usr, _ := user.Current()
	currentDir := fmt.Sprintf("%v/Movies/youtubedr", usr.HomeDir)
	log.Println("download to dir=", currentDir)
	y := NewYoutube(true)
	//arg := flag.Arg(0)
	arg := `https://www.youtube.com/watch?v=UTSqEWQNSCs`
	if err := y.DecodeURL(arg); err != nil {
		fmt.Println("err:", err)
	}
	if err := y.StartDownload(filepath.Join(currentDir, "dl.mp4")); err != nil {
		fmt.Println("err:", err)
	}

}

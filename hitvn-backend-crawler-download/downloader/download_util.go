package main

import (
	"path/filepath"
	"regexp"
	"strings"
)

/**
 * Extracts a file name from url
 */
func GetFileName(url string, videoType string) string{

	if videoType == ".m3u8" {
		if strings.Contains(url, "/d1.vnecdn.net/"){
			var re = regexp.MustCompile(`(?m)[a-z0-9\-\_]+\-\d{8,}`)
			match := re.FindString(url)

			return  match + ".mp4"
		} else if strings.Contains(url, "ss-hls.catscdn.vn") {
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

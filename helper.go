package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func ChanErr(ccc chan int) {
	if ccc != nil {
		ccc <- 1
	}
}

func GetFileContentType(buffer []byte) string {
	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	return contentType
}

func FileCount(dir string) int {
	count := 0
	_ = filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				count += 1
			}
			return nil
		})
	return count
}

func ImageExists(filename string) (bool, int64) {
	log.Debugf("Checking if %s exists...", filename)
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, 0
	}
	log.Debugf("file %s exists!", filename)
	return !info.IsDir(), info.ModTime().Unix()
}

func GenWebpAbs(RawImagePath string, ExhaustPath string,
	ImgFilename string, reqURI string, ModifiedTime int64) (string, string) {
	// get file mod time
	//var ModifiedTime int64
	if ModifiedTime == 0 {
		stat, _ := os.Stat(RawImagePath)
		ModifiedTime = stat.ModTime().Unix()
	}

	// webpFilename: abc.jpg.png -> abc.jpg.png1582558990.webp
	var WebpFilename = fmt.Sprintf("%s.%d.webp", ImgFilename, ModifiedTime)
	cwd, _ := os.Getwd()

	// /home/webp_server/exhaust/path/to/tsuki.jpg.1582558990.webp
	// Custom Exhaust: /path/to/exhaust/web_path/web_to/tsuki.jpg.1582558990.webp
	WebpAbsolutePath := path.Clean(path.Join(ExhaustPath, path.Dir(reqURI), WebpFilename))
	return cwd, WebpAbsolutePath
}

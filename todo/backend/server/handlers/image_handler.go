package handlers

import (
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"
)

type dailyImageHandler struct {
	image        []byte
	modifiedTime time.Time
}

const IMAGE_CACHE_PATH = "cache/image"
const IMAGE_SOURCE_URL = "https://picsum.photos/1200"

func NewDailyImageHandler() *dailyImageHandler {
	d := &dailyImageHandler{}
	stats, statsErr := os.Stat(IMAGE_CACHE_PATH)

	if cacheImageExists(statsErr) && wasCacheImageModifiedToday(stats) {
		err := d.loadImageFromCache(stats)
		if err == nil {
			return d
		}
	}
	d.fetchAndCacheImage()
	return d
}

func (d *dailyImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.refreshImageIfNewDay()
	w.Write(d.image)
}

func cacheImageExists(statsErr error) bool {
	if statsErr != nil {
		log.Println("No cache image found")
	}
	return statsErr == nil
}

func wasCacheImageModifiedToday(stats fs.FileInfo) bool {
	modTime := stats.ModTime()
	currTime := time.Now()

	return sameDate(modTime, currTime)
}

func sameDate(t1 time.Time, t2 time.Time) bool {
	return t1.Day() == t2.Day() && t1.Month() == t2.Month() && t1.Year() == t2.Year()
}

func (d *dailyImageHandler) fetchAndCacheImage() {
	resp, err := http.Get(IMAGE_SOURCE_URL)
	if err != nil {
		log.Printf("failed to fetch image: %v", err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body of image request: %v", err)
		return
	}
	d.image = body
	d.modifiedTime = time.Now()

	cacheImage(body)
}

func cacheImage(image []byte) {
	err := os.WriteFile(IMAGE_CACHE_PATH, image, 0644)
	if err != nil {
		log.Printf("failed to save image to cache: %v", err)
	}
}

func (d *dailyImageHandler) loadImageFromCache(stats fs.FileInfo) error {
	image, err := os.ReadFile(IMAGE_CACHE_PATH)
	if err != nil {
		log.Printf("failed to load image from cache: %v", err)
	} else {
		d.image = image
		d.modifiedTime = stats.ModTime()
	}
	return err
}

func (d *dailyImageHandler) refreshImageIfNewDay() {
	if !sameDate(time.Now(), d.modifiedTime) {
		d.fetchAndCacheImage()
	}
}

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	version    = "dev"
	targetURLs = []string{
		"https://www.cl.cam.ac.uk/",
		"https://www.srcf.net/",
		"https://sobornost.net/",
	}
)

func main() {
	go testLoop()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
}

func testLoop() {
	for {
		offset := float64(120) / float64(len(targetURLs))
		duration := time.Duration(offset * float64(time.Second))
		for _, targetURL := range targetURLs {
			if !testURL(targetURL) {
				fmt.Println("Your clock is skewed compared to " + targetURL)
			}
			time.Sleep(duration)
		}
	}
}

func testURL(targetURL string) bool {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Cache-Control", "no-store, max-age=0")
	req.Header.Set("User-Agent", "timecheckd/"+version+" (+https://github.com/CHTJonas/timecheckd)")

	client := getHTTPClient("")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if dateString := resp.Header.Get("Date"); dateString != "" {
		t := parseHTTPDate(dateString)
		d := time.Now().UTC().Sub(*t)
		fmt.Println("Debug:", targetURL, "time diff is", d)
		if d > 10*time.Second || d < -10*time.Second {
			return false
		}
		return true
	}

	return false
}

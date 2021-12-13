package main

import (
	"context"
	"fmt"
	"net"
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

func parseHTTPDate(dateString string) *time.Time {
	layout := "Mon, 02 Jan 2006 15:04:05 GMT"
	t, err := time.Parse(layout, dateString)
	if err != nil {
		panic(err)
	}
	return &t
}

func getHTTPClient(networkOverride string) *http.Client {
	dialer := &net.Dialer{
		KeepAlive: -1,
	}
	dialCtx := func(ctx context.Context, network, addr string) (net.Conn, error) {
		if networkOverride != "" {
			network = networkOverride
		}
		return dialer.DialContext(ctx, network, addr)
	}
	transport := &http.Transport{
		DisableKeepAlives: true,
		DialContext:       dialCtx,
	}
	return &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Duration(5*time.Second) * time.Millisecond,
	}
}

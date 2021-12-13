package main

import (
	"context"
	"net"
	"net/http"
	"time"
)

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

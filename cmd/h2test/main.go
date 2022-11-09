package main

import (
	"context"
	"fmt"
	"nebula/third_party/http"
)

func main() {
	c := http.NewSimpleClient(context.Background(), http.DefaultClientConfig, http.FingerprintChrome)
	resp, err := c.Builder("https://http3.is").SendAndClose()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Request.URL, resp.Status, resp.Proto)
	resp, err = c.Builder("https://http3.is").SendAndClose()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Request.URL, resp.Status, resp.Proto)
	resp, err = c.Builder("https://www.google.com").SendAndClose()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Request.URL, resp.Status, resp.Proto)
	resp, err = c.Builder("https://www.google.com").SendAndClose()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Request.URL, resp.Status, resp.Proto)
}

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	proxyAddr  = "localhost:8080"
	targetAddr = "echo.test"
)

func main() {
	targetURL, _ := url.Parse("http://" + targetAddr)
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host

		if req.Body != nil && req.Header.Get("Content-Encoding") == "" {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
				return
			}
			_ = req.Body.Close()

			if len(body) > 0 {
				var compressedBody bytes.Buffer
				writer := gzip.NewWriter(&compressedBody)
				_, err = writer.Write(body)
				_ = writer.Close()

				if err != nil {
					log.Printf("Error compressing body: %v", err)
					return
				}

				req.Body = io.NopCloser(&compressedBody)
				req.ContentLength = int64(compressedBody.Len())
				req.Header.Set("Content-Encoding", "gzip")
			} else {
				req.Body = io.NopCloser(bytes.NewReader(body))
			}
		}
	}

	fmt.Printf("Starting proxy server on %s", proxyAddr)
	log.Fatal(http.ListenAndServe(proxyAddr, proxy))
}

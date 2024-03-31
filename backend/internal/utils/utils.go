package utils

import (
	"log"

	"github.com/valyala/fasthttp"
)

func SendPostRequest(url string, jsonData []byte) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	req.SetRequestURI(url)
	req.Header.SetContentType("application/json")
	req.SetBody(jsonData)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		log.Printf("Error sending request: %v\n", err)
		return
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		log.Printf("Error response from server, status code: %d\n", resp.StatusCode())
		return
	}

	log.Println("Response from server received successfully.")
}

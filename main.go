package main

import (
	"github.com/mxmCherry/openrtb/openrtb2"
	"log"
	"net/http"
)
func BidHandler(w http.ResponseWriter, r *http.Request) {
	/*
		JSON transform from request
		1. imp.secure -> ext.is_secure (omit if secure is not defined)
		2. imp[0].id -> ext.id
		3. ext.cb - random string
		4. ext.user-agent - User-Agent header
	*/
	var bidRequest openrtb2.BidRequest
	log.Println(bidRequest)
	// response as json
}

func main() {
	http.HandleFunc("/bid", BidHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
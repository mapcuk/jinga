package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mxmCherry/openrtb/openrtb2"
)

func BidHandler(w http.ResponseWriter, r *http.Request) {

	//validate input

	if contentType := r.Header.Get("Content-Type"); contentType != "" && contentType != "application/json" {
		msg := "Content-Type header is not application/json"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		log.Printf("Invalid content type %s\n", contentType)
		return
	}

	if r.Body == nil {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		log.Printf("Nil request body")
		return
	}

	var bidRequest openrtb2.BidRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&bidRequest); err != nil {
		http.Error(w, "Malformed content", http.StatusBadRequest)
		log.Printf("Malformed content %+v, %v\n", bidRequest, err)
		return
	}
	//not needed but just in case
	defer r.Body.Close()

	if len(bidRequest.Imp) == 0 {
		http.Error(w, "Malformed content", http.StatusBadRequest)
		log.Printf("Malformed content %+v\n", bidRequest)
		return
	}

	//gather needed data

	isSecure := bidRequest.Imp[0].Secure
	id := bidRequest.Imp[0].ID
	cb := String(10)
	userAgent := r.UserAgent()

	t := ext{
		ID:        id,
		Cb:        cb,
		IsSecure:  isSecure,
		UserAgent: userAgent,
	}

	tBytes, err := json.Marshal(t)
	if err != nil {
		http.Error(w, "Please, try again later", http.StatusInternalServerError)
		log.Printf("Error marshalling to type ext, %v", err)
		return
	}

	bidRequest.Ext = json.RawMessage(tBytes)

	//construct response
	res, err := json.Marshal(bidRequest)
	if err != nil {
		http.Error(w, "Please, try again later", http.StatusInternalServerError)
		log.Printf("Error marshalling response, %v", err)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		log.Printf("Error writing response, %v", err)
	}
}

type ext struct {
	ID        string `json:"id"`
	Cb        string `json:"cb"`
	IsSecure  *int8  `json:"is_secure,omitempty"`
	UserAgent string `json:"user-agent"`
}

func main() {
	http.HandleFunc("/bid", BidHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

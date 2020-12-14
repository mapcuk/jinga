package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestBidHandler_NilBody(t *testing.T) {
	req, err := http.NewRequest("GET", "/bid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BidHandler)

	handler.ServeHTTP(rr, req)

	// must return 404 for nil body
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestBidHandler_WrongContentType(t *testing.T) {
	reader, _ := os.Open("input.json")
	req, err := http.NewRequest("POST", "/bid", reader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "text/plain")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BidHandler)

	handler.ServeHTTP(rr, req)

	// must return 415 for wrong content-type
	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnsupportedMediaType)
	}
}

func TestBidHandler_EmptyImpl(t *testing.T) {
	reader, _ := os.Open("input_empty_impl.json")
	req, err := http.NewRequest("POST", "/bid", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BidHandler)

	handler.ServeHTTP(rr, req)

	// must return 415 for wrong content-type
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestBidHandler_AllIsOK(t *testing.T) {
	reader, _ := os.Open("input.json")
	req, err := http.NewRequest("POST", "/bid", reader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BidHandler)

	handler.ServeHTTP(rr, req)

	// must return 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `\{"id":"42","imp":\[\{"id":"someId","banner":\{"w":1280,"h":720\},"secure":1\}\],"ext":\{"id":"someId","cb":".*?","is_secure":true,"user-agent":"Mozilla/5.0 \(X11; Ubuntu; Linux x86_64; rv:83.0\) Gecko/20100101 Firefox/83.0"\}\}`
	matched, _ := regexp.MatchString(expected, rr.Body.String())
	if !matched {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestBidHandler_MissingSecure(t *testing.T) {
	reader, _ := os.Open("input_missing_secure.json")
	req, err := http.NewRequest("POST", "/bid", reader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BidHandler)

	handler.ServeHTTP(rr, req)

	// must return 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if strings.Contains(rr.Body.String(), "is_secure") {
		t.Error("result is not supposed to container is_secure field")
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

var router *mux.Router

func TestMain(m *testing.M) {
	fmt.Println("--- Start Tests")
	setup()

	// run the test cases defined in this file
	retCode := m.Run()

	tearDown()

	// call with result of m.Run()
	os.Exit(retCode)
}

func setup() {
	fmt.Println("--- --- setup")
	router = makeRouter()
}

func tearDown() {
	fmt.Println("--- --- tear down")
}

func buildRequest(method, endpoint string, payload io.Reader, t *testing.T) *http.Request {
	req, err := http.NewRequest(method, endpoint, payload)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}

	return req
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}
func TestGetAppPage(t *testing.T) {
	fmt.Println("start TestGetAppReviewsOfClass")
	var method = "GET"
	var endpoint = "/hitec/crawl/app-page/google-play/%s"

	/*
	 * test for success CHECK 1
	 */
	endpointCheckOne := fmt.Sprintf(endpoint, "com.whatsapp")
	req := buildRequest(method, endpointCheckOne, nil, t)
	rr := executeRequest(req)

	//Confirm the response has the right status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	var appPages AppPage
	err := json.NewDecoder(rr.Body).Decode(&appPages)
	if err != nil {
		t.Errorf("Did not receive a proper formed json")
	}
}

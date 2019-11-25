package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/OlegSchmidt/soup"
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
	fmt.Println("start TestGetAppPage")
	var method = "GET"
	var endpoint = "/hitec/crawl/app-page/google-play/%s"

	/*
	 * test for success CHECK 1
	 */
	for _, appId := range []string{
		"com.whatsapp",             // free
		"com.ustwo.monumentvalley", // paid, contains in-app-purchases
		"com.does.not.exists.122",  // does not exist
	} {
		endpointCheckOne := fmt.Sprintf(endpoint, appId)
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
}

var mailformedHTML = `
<html>
  <head>
    <title>Sample "Hello, World" Application</title>
  </head>
<html>
`

func TestGetRating(t *testing.T) {
	for _, document := range []string{
		mailformedHTML,
	} {
		_, error := getRating(soup.HTMLParse(document))
		if error == nil {
			t.Errorf("should contain error")
		}
	}
}
func TestGetAppName(t *testing.T) {
	for _, document := range []string{
		mailformedHTML,
	} {
		_, error := getAppName(soup.HTMLParse(document))
		if error == nil {
			t.Errorf("should contain error")
		}
	}
}
func TestGetCategory(t *testing.T) {
	for _, document := range []string{
		mailformedHTML,
	} {
		_, error := getCategory(soup.HTMLParse(document))
		if error == nil {
			t.Errorf("should contain error")
		}
	}
}
func TestGetUsk(t *testing.T) {
	for _, document := range []string{
		mailformedHTML,
	} {
		_, error := getUsk(soup.HTMLParse(document))
		if error == nil {
			t.Errorf("should contain error")
		}
	}
}
func TestGetDescription(t *testing.T) {
	for _, document := range []string{
		mailformedHTML,
	} {
		_, error := getDescription(soup.HTMLParse(document))
		if error == nil {
			t.Errorf("should contain error")
		}
	}
}
func TestGetWhatsNew(t *testing.T) {
	for _, document := range []string{
		mailformedHTML,
	} {
		_, error := getWhatsNew(soup.HTMLParse(document))
		if error == nil {
			t.Errorf("should contain error")
		}
	}
}
func TestGetStarsCount(t *testing.T) {
	for _, document := range []string{
		mailformedHTML,
	} {
		_, error := getStarsCount(soup.HTMLParse(document))
		if error == nil {
			t.Errorf("should contain error")
		}
	}
}
func TestCountPerRating(t *testing.T) {
	for _, document := range []string{
		mailformedHTML,
	} {
		_, error := getCountPerRating(soup.HTMLParse(document))
		if error == nil {
			t.Errorf("should contain error")
		}
	}
}
func TestDeveloperName(t *testing.T) {
	for _, document := range []string{
		mailformedHTML,
	} {
		_, error := getCountPerRating(soup.HTMLParse(document))
		if error == nil {
			t.Errorf("should contain error")
		}
	}
}
func TestTopDeveloper(t *testing.T) {
	for _, document := range []string{
		mailformedHTML,
	} {
		_, error := getTopDeveloper(soup.HTMLParse(document))
		if error == nil {
			t.Errorf("should contain error")
		}
	}
}

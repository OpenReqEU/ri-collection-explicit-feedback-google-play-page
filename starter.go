package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	requestError = "The request could not be recovered"
)

func main() {
	log.Fatal(http.ListenAndServe(":9622", makeRouter()))
}

func makeRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/hitec/crawl/app-page/google-play/{package_name}", getAppPage).Methods("GET")
	return router
}

func recoverAPICall(w http.ResponseWriter, page AppPage) {
	if r := recover(); r != nil {
		w.WriteHeader(http.StatusInternalServerError)
		page.Errors = append(page.Errors, requestError)
	}
}

func getAppPage(w http.ResponseWriter, r *http.Request) {
	appPage := AppPage{}
	w.Header().Set("Content-Type", "application/json")
	defer recoverAPICall(w, appPage)

	// get request param
	params := mux.Vars(r)
	packageName := params["package_name"]

	// crawl app reviews
	appPage = Crawl(packageName)
	serveResponse(w, appPage, http.StatusOK)
}

// serves the generated content
func serveResponse(writer http.ResponseWriter, page AppPage, status int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(page)
}

package main

import (
	"log"
	"net/http"

	"encoding/json"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	log.SetOutput(os.Stdout)

	router := mux.NewRouter()
	router.HandleFunc("/hitec/crawl/app-page/google-play/{package_name}", getAppPage).Methods("GET")

	log.Fatal(http.ListenAndServe(":9622", router))
}

func getAppPage(w http.ResponseWriter, r *http.Request) {
	// get request param
	params := mux.Vars(r)
	packageName := params["package_name"]

	// crawl app reviews
	appPage := Crawl(packageName)

	// write the response
	w.Header().Set("Content-Type", "application/json")
	if appPage.Description != "" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(appPage)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(`{"message": "could not retrieve app page"}`)
	}
}

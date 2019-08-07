package main

import (
	"encoding/json"
	"log"
	"net/http"

	"os"

	"github.com/gorilla/mux"
)

func main() {
	parameters := os.Args
	if len(parameters) >= 2 && parameters[len(os.Args)-2] == "local" {
		Crawl(parameters[len(os.Args)-1], true)
	} else {
		log.SetOutput(os.Stdout)
		router := mux.NewRouter()
		router.HandleFunc("/hitec/crawl/app-page/google-play/{package_name}", getAppPage).Methods("GET")
		log.Fatal(http.ListenAndServe(":9622", router))
	}
}

func recoverAPICall(w http.ResponseWriter) {
	if r := recover(); r != nil {
		log.Println("recovered from ", r)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(`{"message": "could not retrieve app page"}`)
	}
}

func getAppPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recoverAPICall(w)

	// get request param
	params := mux.Vars(r)
	packageName := params["package_name"]

	// crawl app reviews
	appPage := Crawl(packageName)
	if appPage.Description != "" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(appPage)
	}
}

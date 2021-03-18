package main

import (
	"log"
	"main/api/cases"
	"main/api/policy"
	"net/http"
	"os"
)

// Main program.
func main() {
	//get port
	port := os.Getenv("PORT")
	//branch if there isn't a port and set it to 8080
	if port == "" {
		port = "8080"
	}
	//declare structures with handlers
	var cases cases.Cases
	var policy policy.Policy
	//handle corona cases
	http.HandleFunc("/corona/v1/country/", cases.Handler)
	http.HandleFunc("/corona/v1/policy/", policy.Handler)
	//ends program if it can't open port
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

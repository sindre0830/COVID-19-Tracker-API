package main

import (
	"log"
	"main/api/cases"
	"main/api/diag"
	"main/api/policy"
	"net/http"
	"os"
	"time"
)

// init runs once at startup.
func init() {
	//set varible to current time (for uptime)
	diag.StartTime = time.Now()
}
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
	var diagnosis diag.Diagnosis
	//handle corona cases
	http.HandleFunc("/corona/v1/country/", cases.Handler)
	//handle corona policy
	http.HandleFunc("/corona/v1/policy/", policy.Handler)
	//handle program diagnosis
	http.HandleFunc("/corona/v1/diag/", diagnosis.Handler)
	//ends program if it can't open port
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

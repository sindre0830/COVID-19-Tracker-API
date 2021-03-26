package main

import (
	"log"
	"main/api/cases"
	"main/api/diag"
	"main/api/notification"
	"main/api/policy"
	"main/dict"
	"net/http"
	"os"
	"time"
)

// init runs once at startup.
func init() {
	//set varible to current time (for uptime)
	diag.StartTime = time.Now()
	//setup connection with firebase and branch if an error occured
	err := notification.DB.Setup()
	if err != nil {
		defer notification.DB.Client.Close()
		log.Fatalln(err)
	}
}

// Main program.
func main() {
	//get port and branch if there isn't a port and set it to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	//
	dict.URL = dict.URL + ":" + port
	//schedule checks every second for possible webhooks to execute
	go notification.Schedule()
	//handle corona cases
	http.HandleFunc("/corona/v1/country/", cases.MethodHandler)
	//handle corona policy
	http.HandleFunc("/corona/v1/policy/", policy.MethodHandler)
	//handle program diagnosis
	http.HandleFunc("/corona/v1/diag/", diag.MethodHandler)
	//handle webhook methods
	http.HandleFunc("/corona/v1/notifications/", notification.MethodHandler)
	//ends program if it can't open port
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

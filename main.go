package main

import (
	"log"
	"main/api/cases"
	"main/api/diag"
	"main/api/notification"
	"main/api/policy"
	"main/firebase"
	"net/http"
	"os"
	"time"
)

// init runs once at startup.
func init() {
	//set varible to current time (for uptime)
	diag.StartTime = time.Now()
	notification.Secret = []byte{1, 2, 3, 4, 5} // not a good secret!
	//setup connection with firebase and branch if an error occured
	err := firebase.DB.Setup()
	if err != nil {
		defer firebase.DB.Client.Close()
		log.Fatalln(err)
	}
}
// Main program.
func main() {
	//get port
	port := os.Getenv("PORT")
	//branch if there isn't a port and set it to 8080
	if port == "" {
		port = "8080"
	}
	//handle corona cases
	http.HandleFunc("/corona/v1/country/", cases.MethodHandler)
	//handle corona policy
	http.HandleFunc("/corona/v1/policy/", policy.MethodHandler)
	//handle program diagnosis
	http.HandleFunc("/corona/v1/diag/", diag.MethodHandler)
	//handle webhook methods
	http.HandleFunc("/corona/v1/notifications/", notification.MethodHandler)
	//handle webhook methods
	http.HandleFunc("/corona/v1/service/", notification.ServiceHandler)
	//ends program if it can't open port
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

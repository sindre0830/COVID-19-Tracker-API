package main

import (
	"context"
	"log"
	"main/api/cases"
	"main/api/diag"
	"main/api/notification"
	"main/api/policy"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// init runs once at startup.
func init() {
	//set varible to current time (for uptime)
	diag.StartTime = time.Now()
	notification.Secret = []byte{1, 2, 3, 4, 5} // not a good secret!
}
// Main program.
func main() {
	//get port
	port := os.Getenv("PORT")
	//branch if there isn't a port and set it to 8080
	if port == "" {
		port = "8080"
	}
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
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

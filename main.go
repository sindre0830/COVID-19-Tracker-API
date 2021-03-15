package main

import (
	"log"
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
	//ends program if it can't open port
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

package main

import (
	"log"
	"main/api"
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

	var cases api.Cases
	cases.Handler("Norway", "2020-12-01", "2021-01-31")

	//ends program if it can't open port
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

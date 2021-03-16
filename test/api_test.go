package api_test

import (
	"fmt"
	"main/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_casesHandler(t *testing.T) {
	var cases api.Cases
	//store expected data to check against
	data := map[string]int {
		//test path
		"http://localhost:8080/corona/v1/country/": http.StatusBadRequest,
		"http://localhost:8080/corona/v1/country/norway/": http.StatusBadRequest,
		"http://localhost:8080/corona/v1/country/norway": http.StatusOK,
		//test country edge case
		"http://localhost:8080/corona/v1/country/italy": http.StatusOK,
		"http://localhost:8080/corona/v1/country/NORWAY": http.StatusOK,
		"http://localhost:8080/corona/v1/country/nor": http.StatusInternalServerError,
		"http://localhost:8080/corona/v1/country/usa": http.StatusInternalServerError,
		//test parameters
		"http://localhost:8080/corona/v1/country/norway?": http.StatusOK,
		"http://localhost:8080/corona/v1/country/norway?abc=something": http.StatusBadRequest,
		"http://localhost:8080/corona/v1/country/norway?scope=2020-01-01-2021-01-01": http.StatusOK,
		"http://localhost:8080/corona/v1/country/norway?scope=2020-01-01-2021-01-01?abc=something": http.StatusBadRequest,
	}
	//iterate through map and check each key to expected element
	for url, expectedStatus := range data {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating HTTP request in Test_casesHandler")
			return
		}
		recorder := httptest.NewRecorder()
		cases.Handler(recorder, req)
		//branch if we get an unexpected answer
		if recorder.Code != expectedStatus {
			t.Errorf("Expected '%v' but got '%v'. Tested: %v", expectedStatus, recorder.Code, url)
		}
	}
}

package api_test

import (
	"fmt"
	"main/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_PolicyHistory_Get(t *testing.T) {
	//store expected data to check against
	data := map[[3]string]int {
		{"nor", "2021-01-01", "2021-03-01"}: http.StatusOK,
		{"nor", "2022-01-01", "2022-03-01"}: http.StatusNotFound,
		{"nor", "2016-01-01", "2016-03-01"}: http.StatusNotFound,
	}
	//iterate through map and check each key to expected element
	for arrTestData, expectedStatus := range data {
		var policyHistory api.PolicyHistory
		_, status, _ := policyHistory.Get(arrTestData[0], arrTestData[1], arrTestData[2])
		//branch if we get an unexpected answer
		if status != expectedStatus {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'.", expectedStatus, status, arrTestData)
		}
	}
}

func Test_PolicyCurrent_Get(t *testing.T) {
	//store expected data to check against
	data := map[string]int {
		"": http.StatusNotFound,
		"norway": http.StatusNotFound,
		"nor": http.StatusOK,
	}
	//iterate through map and check each key to expected element
	for country, expectedStatus := range data {
		var policyCurrent api.PolicyCurrent
		_, status, _ := policyCurrent.Get(country)
		//branch if we get an unexpected answer
		if status != expectedStatus {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'", expectedStatus, status, country)
		}
	}
}

func Test_CountryNameDetails_Get(t *testing.T) {
	//store expected data to check against
	data := map[string]int {
		"": http.StatusNotFound,
		"norwayyy": http.StatusNotFound,
		"norway": http.StatusOK,
	}
	//iterate through map and check each key to expected element
	for country, expectedStatus := range data {
		var countryNameDetails api.CountryNameDetails
		status, _ := countryNameDetails.Get(country)
		//branch if we get an unexpected answer
		if status != expectedStatus {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'", expectedStatus, status, country)
		}
	}
}

func Test_casesHandler(t *testing.T) {
	//store expected data to check against
	data := map[string]int {
		//test path
		"http://localhost:8080/corona/v1/country/": http.StatusBadRequest,
		"http://localhost:8080/corona/v1/country/norway/": http.StatusBadRequest,
		"http://localhost:8080/corona/v1/country/norway": http.StatusOK,
		//test country edge case
		"http://localhost:8080/corona/v1/country/italy": http.StatusOK,
		"http://localhost:8080/corona/v1/country/NORWAY": http.StatusOK,
		"http://localhost:8080/corona/v1/country/nor": http.StatusBadRequest,
		"http://localhost:8080/corona/v1/country/usa": http.StatusBadRequest,
		//test parameters
		"http://localhost:8080/corona/v1/country/norway?": http.StatusOK,
		"http://localhost:8080/corona/v1/country/norway?abc=something": http.StatusBadRequest,
		"http://localhost:8080/corona/v1/country/norway?scope=2020-01-01-2021-01-01": http.StatusOK,
		"http://localhost:8080/corona/v1/country/norway?scope=2020-01-01-2021-01-01?abc=something": http.StatusBadRequest,
	}
	//iterate through map and check each key to expected element
	for url, expectedStatus := range data {
		var cases api.Cases
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

package cases_test

import (
	"fmt"
	"main/api/cases"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Cases_Handler(t *testing.T) {
	//store expected data to check against
	data := map[string]int{
		//test path
		"http://localhost:8080/corona/v1/country/":        http.StatusNotFound,
		"http://localhost:8080/corona/v1/country/norway/": http.StatusBadRequest,
		"http://localhost:8080/corona/v1/country/norway":  http.StatusOK,
		//test country edge case
		"http://localhost:8080/corona/v1/country/italy":  http.StatusOK,
		"http://localhost:8080/corona/v1/country/NORWAY": http.StatusOK,
		"http://localhost:8080/corona/v1/country/nor":    http.StatusBadRequest,
		"http://localhost:8080/corona/v1/country/usa":    http.StatusOK,
		//test parameters
		"http://localhost:8080/corona/v1/country/norway?":                                          http.StatusOK,
		"http://localhost:8080/corona/v1/country/norway?abc=something":                             http.StatusBadRequest,
		"http://localhost:8080/corona/v1/country/norway?scope=2020-01-01-2021-01-01":               http.StatusOK,
		"http://localhost:8080/corona/v1/country/norway?scope=2020-01-01-2021-01-01?abc=something": http.StatusBadRequest,
	}
	//iterate through map and check each key to expected element
	for url, expectedStatus := range data {
		var cases cases.Cases
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

func Test_CasesHistory_Get(t *testing.T) {
	//store expected data to check against
	data := map[[3]string]int{
		{"Norway", "2021-01-01", "2021-03-01"}: http.StatusOK,
		{"Norway", "2022-01-01", "2022-03-01"}: http.StatusOK,
		{"Norway", "2016-01-01", "2016-03-01"}: http.StatusOK,
		{"norway", "2021-01-01", "2021-03-01"}: http.StatusBadRequest,
		{"USA", "2021-01-01", "2021-03-01"}: http.StatusBadRequest,
		{"US", "2021-01-01", "2021-03-01"}: http.StatusOK,
	}
	//iterate through map and check each key to expected element
	for arrTestData, expectedStatus := range data {
		var casesHistory cases.CasesHistory
		_, _, status, _ := casesHistory.Get(arrTestData[0], arrTestData[1], arrTestData[2])
		//branch if we get an unexpected answer
		if status != expectedStatus {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'.", expectedStatus, status, arrTestData)
		}
	}
}

func Test_CasesTotal_Get(t *testing.T) {
	//store expected data to check against
	data := map[string]int{
		"Norway": http.StatusOK,
		"norway": http.StatusBadRequest,
		"nor": http.StatusBadRequest,
		"USA": http.StatusBadRequest,
		"US": http.StatusOK,
		"Us": http.StatusBadRequest,
	}
	//iterate through map and check each key to expected element
	for testData, expectedStatus := range data {
		var casesTotal cases.CasesTotal
		status, _ := casesTotal.Get(testData)
		//branch if we get an unexpected answer
		if status != expectedStatus {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'.", expectedStatus, status, testData)
		}
	}
}

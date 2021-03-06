package policy_test

import (
	"fmt"
	"main/api/policy"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Policy_Handler(t *testing.T) {
	//store expected data to check against
	testData := map[string]int{
		//test path
		"http://localhost:8080/corona/v1/policy/":        http.StatusNotFound,
		"http://localhost:8080/corona/v1/policy/norway/": http.StatusBadRequest,
		"http://localhost:8080/corona/v1/policy/norway":  http.StatusOK,
		//test country edge case
		"http://localhost:8080/corona/v1/policy/italy":  http.StatusOK,
		"http://localhost:8080/corona/v1/policy/nor":    http.StatusNotFound,
		"http://localhost:8080/corona/v1/policy/usa":    http.StatusOK,
		//test parameters
		"http://localhost:8080/corona/v1/policy/norway?":                                          http.StatusOK,
		"http://localhost:8080/corona/v1/policy/norway?abc=something":                             http.StatusBadRequest,
		"http://localhost:8080/corona/v1/policy/norway?scope=2020-01-01-2021-01-01":               http.StatusOK,
		"http://localhost:8080/corona/v1/policy/norway?scope=2020-01-01-2021-01-01?abc=something": http.StatusBadRequest,
	}
	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {
		var policy policy.Policy
		req, err := http.NewRequest("GET", test, nil)
		if err != nil {
			fmt.Println("Error creating HTTP request in Test_Policy_Handler")
			return
		}
		recorder := httptest.NewRecorder()
		policy.Handler(recorder, req)
		//branch if we get an unexpected answer that is not timed out
		if recorder.Code != expectedStatus && recorder.Code != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: %v", expectedStatus, recorder.Code, test)
		}
	}
}

func Test_PolicyHistory_Get(t *testing.T) {
	//store expected data to check against
	testData := map[[2]string]int{
		{"2021-01-01", "2021-03-01"}: http.StatusOK,
		{"2022-01-01", "2022-03-01"}: http.StatusOK,
		{"2016-01-01", "2016-03-01"}: http.StatusOK,
	}
	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {
		var policyHistory policy.PolicyHistory
		status, _ := policyHistory.Get(test[0], test[1])
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'.", expectedStatus, status, test)
		}
	}
}

func Test_PolicyCurrent_Get(t *testing.T) {
	//store expected data to check against
	testData := map[[2]string]int{
		{"", "2021-01-01"}:       http.StatusNotFound,
		{"norway", "2021-01-01"}: http.StatusOK,
		{"nor", "2021-01-01"}:    http.StatusOK,
	}
	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {
		var policyCurrent policy.PolicyCurrent
		status, _ := policyCurrent.Get(test[0], test[1])
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'", expectedStatus, status, test)
		}
	}
}

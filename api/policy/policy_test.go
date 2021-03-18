package policy_test

import (
	"main/api/policy"
	"net/http"
	"testing"
)

func Test_PolicyHistory_Get(t *testing.T) {
	//store expected data to check against
	data := map[[3]string]int{
		{"2021-01-01", "2021-03-01"}: http.StatusOK,
		{"2022-01-01", "2022-03-01"}: http.StatusNotFound,
		{"2016-01-01", "2016-03-01"}: http.StatusNotFound,
	}
	//iterate through map and check each key to expected element
	for arrTestData, expectedStatus := range data {
		var policyHistory policy.PolicyHistory
		status, _ := policyHistory.Get(arrTestData[0], arrTestData[1])
		//branch if we get an unexpected answer
		if status != expectedStatus {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'.", expectedStatus, status, arrTestData)
		}
	}
}

func Test_PolicyCurrent_Get(t *testing.T) {
	//store expected data to check against
	data := map[[2]string]int{
		{"", "2021-01-01"}:       http.StatusNotFound,
		{"norway", "2021-01-01"}: http.StatusNotFound,
		{"nor", "2021-01-01"}:    http.StatusOK,
	}
	//iterate through map and check each key to expected element
	for arrTestData, expectedStatus := range data {
		var policyCurrent policy.PolicyCurrent
		status, _ := policyCurrent.Get(arrTestData[0], arrTestData[1])
		//branch if we get an unexpected answer
		if status != expectedStatus {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'", expectedStatus, status, arrTestData)
		}
	}
}

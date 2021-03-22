package countryinfo_test

import (
	"main/api/countryinfo"
	"net/http"
	"testing"
)

func Test_CountryNameDetails_Get(t *testing.T) {
	//store expected data to check against
	data := map[string]int {
		"": http.StatusNotFound,
		"norwayyy": http.StatusNotFound,
		"norway": http.StatusOK,
	}
	//iterate through map and check each key to expected element
	for country, expectedStatus := range data {
		var countryNameDetails countryinfo.CountryNameDetails
		status, _ := countryNameDetails.Get(country)
		//branch if we get an unexpected answer
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'", expectedStatus, status, country)
		}
	}
}

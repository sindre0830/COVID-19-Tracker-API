package fun_test

import (
	"main/fun"
	"net/url"
	"testing"
)

func Test_LimitDecimals(t *testing.T) {
	//store expected data to check against
	testData := map[float64]float64 {
		0.0: 		0.0,
		54.123453:  54.12,
		-323.1111: 	-323.11,
	}
	//iterate through map and check each key to expected element
	for test, expectedResult := range testData {
		result := fun.LimitDecimals(test)
		//branch if we get an unexpected answer that is not timed out
		if result != expectedResult {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'", expectedResult, result, test)
		}
	}
}

func Test_ParseURL(t *testing.T) {
	//store expected data to check against
	testData := map[string]bool {
		"http://localhost:8080/corona/v1/country": false,
		"http://localhost:8080/corona/v1/country/norway": true,
		"http://localhost:8080/corona/v1/country/norway/": false,
		"http://localhost:8080/corona/v1/country/norway?scope=date": true,
		"http://localhost:8080/corona/v1/country/norway?rand=idk": false,
	}
	//iterate through map and check each key to expected element
	for test, expectedResult := range testData {
		parsedURL, _ := url.ParseRequestURI(test)
		_, _, _, result := fun.ParseURL(parsedURL)
		//branch if we get an unexpected answer that is not timed out
		if result != nil && expectedResult != false || result == nil && expectedResult != true {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'", expectedResult, result, test)
		}
	}
}

func Test_ValidateCountry(t *testing.T) {
	//store expected data to check against
	testData := map[string]bool {
		"": 		false,
		"Norway":   true,
	}
	//iterate through map and check each key to expected element
	for test, expectedResult := range testData {
		result := fun.ValidateCountry(test)
		//branch if we get an unexpected answer that is not timed out
		if result != nil && expectedResult != false || result == nil && expectedResult != true {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'", expectedResult, result, test)
		}
	}
}

func Test_ValidateDates(t *testing.T) {
	//store expected data to check against
	testData := map[string]bool {
		"2020-01-01-2020-02-01": true,
		"2020-01-01-2020-02":    false,
		"01-01-2020-2020-02-01": false,
		"2020-01-01-01-02-2020": false,
		"2020-01-01-2O20-02-01": false,
		"2020-01-01-2020-02-00": false,
		"2020-02-01-2020-01-01": false,
	}
	//iterate through map and check each key to expected element
	for test, expectedResult := range testData {
		result := fun.ValidateDates(test)
		//branch if we get an unexpected answer
		if result != nil && expectedResult != false || result == nil && expectedResult != true {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'", expectedResult, result, test)
		}
	}
}

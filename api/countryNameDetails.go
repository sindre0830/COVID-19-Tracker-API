package api

import (
	"encoding/json"
	"net/http"
)

// get will update CountryNameDetails based on input.
func (countryNameDetails *CountryNameDetails) Get(country string) (int, error) {
	url := "https://restcountries.eu/rest/v2/name/" + country + "?fields=name;alpha3Code"
	//gets json output from API and branch if an error occurred
	status, err := countryNameDetails.req(url)
	if err != nil {
		return status, err
	}
	return http.StatusOK, nil
}
// req will request from API based on URL.
func (countryNameDetails *CountryNameDetails) req(url string) (int, error) {
	//gets raw data from API and branch if an error occurred
	data, status, err := requestData(url)
	if err != nil {
		return status, err
	}
	//convert raw data to JSON and branch if an error occurred
	err = json.Unmarshal(data, &countryNameDetails)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

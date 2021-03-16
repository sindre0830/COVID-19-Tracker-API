package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

// get will update object based on input.
func (object *casesHistory) get(country string, startDate string, endDate string) (int, int, int, error) {
	url := "https://covid-api.mmediagroup.fr/v1/history?country=" + country + "&status=Confirmed"
	//gets JSON output from API based on confirmed cases and branch if an error occurred
	status, err := object.req(url)
	if err != nil {
		return 0, 0, status, err
	}
	//branch if output from API is empty and return error
	if object.isEmpty() {
		err = errors.New("object validation: object is empty")
		return 0, 0, http.StatusBadRequest, err
	}
	//get cases between start and end date
	confirmed := object.addCases(startDate, endDate)
	url = "https://covid-api.mmediagroup.fr/v1/history?country=" + country + "&status=Recovered"
	//gets json output from API based on recovered cases and branch if an error occurred
	status, err = object.req(url)
	if err != nil {
		return 0, 0, status, err
	}
	//get cases between start and end date
	recovered := object.addCases(startDate, endDate)
	return confirmed, recovered, http.StatusOK, nil
}
// req will request from API based on URL.
func (object *casesHistory) req(url string) (int, error) {
	//gets raw data from API and branch if an error occurred
	output, status, err := requestData(url)
	if err != nil {
		return status, err
	}
	//convert raw data to JSON and branch if an error occurred
	err = json.Unmarshal(output, &object)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
// isEmpty checks if object is empty.
func (object *casesHistory) isEmpty() bool {
    return object.All.Country == ""
}
// addCases gets cases between two dates.
func (object *casesHistory) addCases(startDate string, endDate string) int {
	n:= object.All.Dates[endDate] - object.All.Dates[startDate]
	//branch if number is negative and convert to absolute
	if n < 0 {
		n *= (-1)
	}
	return n
}

package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

// get will update casesTotal based on input.
func (casesTot *casesTotal) get(country string) (int, error) {
	url := "https://covid-api.mmediagroup.fr/v1/cases?country=" + country
	//gets json output from API and branch if an error occurred
	status, err := casesTot.req(url)
	if err != nil {
		return status, err
	}
	//branch if output from API is empty and return error
	if casesTot.isEmpty() {
		err = errors.New("casesTot validation: casesTot is empty")
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
// req will request from API based on URL.
func (casesTot *casesTotal) req(url string) (int, error) {
	//gets raw data from API and branch if an error occurred
	data, status, err := requestData(url)
	if err != nil {
		return status, err
	}
	//convert raw data to JSON and branch if an error occurred
	err = json.Unmarshal(data, &casesTot)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
// isEmpty checks if casesTotal is empty.
func (casesTot *casesTotal) isEmpty() bool {
    return casesTot.All.Country == ""
}

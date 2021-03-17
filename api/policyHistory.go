package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

// get will update PolicyHistory based on input.
func (policyHistory *PolicyHistory) Get(country string, startDate string, endDate string) (int, error) {
	url := "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/" + startDate + "/" + endDate
	//gets json output from API and branch if an error occurred
	status, err := policyHistory.req(url)
	if err != nil {
		return status, err
	}
	//branch if output from API is empty and return error
	if policyHistory.isEmpty() {
		err = errors.New("object validation: PolicyHistory is empty")
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}
// req will request from API based on URL.
func (policyHistory *PolicyHistory) req(url string) (int, error) {
	//gets raw data from API and branch if an error occurred
	data, status, err := requestData(url)
	if err != nil {
		return status, err
	}
	//convert raw data to JSON and branch if an error occurred
	err = json.Unmarshal(data, &policyHistory)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
// isEmpty checks if PolicyHistory is empty.
func (policyHistory *PolicyHistory) isEmpty() bool {
    return policyHistory.Scale == nil
}

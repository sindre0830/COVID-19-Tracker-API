package policy

import (
	"encoding/json"
	"errors"
	"main/api"
	"net/http"
)

// PolicyHistory stores data about COVID policies for all countries between two dates.
//
// Functionality: Get, req, isEmpty, decreaseDate
type PolicyHistory struct {
	Scale     map[string]map[string]int `json:"scale"`
	Countries []string                  `json:"countries"`
	Data      map[string]map[string]struct {
		DateValue            string  `json:"date_value"`
		CountryCode          string  `json:"country_code"`
		Confirmed            int     `json:"confirmed"`
		Deaths               int     `json:"deaths"`
		StringencyActual     float64 `json:"stringency_actual"`
		Stringency           float64 `json:"stringency"`
		StringencyLegacy     float64 `json:"stringency_legacy"`
		StringencyLegacyDisp float64 `json:"stringency_legacy_disp"`
	} `json:"data"`
}
// get will update PolicyHistory based on input.
func (policyHistory *PolicyHistory) Get(startDate string, endDate string) (int, error) {
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
	data, status, err := api.RequestData(url)
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

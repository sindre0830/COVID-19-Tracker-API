package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

// get will update casesHistory based on input.
func (casesHis *casesHistory) get(country string, startDate string, endDate string) (int, int, int, error) {
	url := "https://covid-api.mmediagroup.fr/v1/history?country=" + country + "&status=Confirmed"
	//gets JSON output from API based on confirmed cases and branch if an error occurred
	status, err := casesHis.req(url)
	if err != nil {
		return 0, 0, status, err
	}
	//branch if output from API is empty and return error
	if casesHis.isEmpty() {
		err = errors.New("casesHis validation: casesHis is empty")
		return 0, 0, http.StatusBadRequest, err
	}
	//get cases between start and end date
	confirmed := casesHis.addCases(startDate, endDate)
	url = "https://covid-api.mmediagroup.fr/v1/history?country=" + country + "&status=Recovered"
	//gets json output from API based on recovered cases and branch if an error occurred
	status, err = casesHis.req(url)
	if err != nil {
		return 0, 0, status, err
	}
	//get cases between start and end date
	recovered := casesHis.addCases(startDate, endDate)
	return confirmed, recovered, http.StatusOK, nil
}
// req will request from API based on URL.
func (casesHis *casesHistory) req(url string) (int, error) {
	//gets raw data from API and branch if an error occurred
	output, status, err := requestData(url)
	if err != nil {
		return status, err
	}
	//convert raw data to JSON and branch if an error occurred
	err = json.Unmarshal(output, &casesHis)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
// isEmpty checks if casesHistory is empty.
func (casesHis *casesHistory) isEmpty() bool {
    return casesHis.All.Country == ""
}
// addCases gets cases between two dates.
func (casesHis *casesHistory) addCases(startDate string, endDate string) int {
	n:= casesHis.All.Dates[endDate] - casesHis.All.Dates[startDate]
	//branch if number is negative and convert to absolute
	if n < 0 {
		n *= (-1)
	}
	return n
}

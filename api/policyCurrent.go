package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// get will update PolicyCurrent based on input.
func (policyCurrent *PolicyCurrent) Get(country string) (time.Time, int, error) {
	//get current time in YYYY-MM-DD format
	currentTime := time.Now()
	pastTime := currentTime.AddDate(0, 0, -10)
	date := pastTime.Format("2006-01-02")
	fmt.Println(date)
	url := "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/" + country + "/" + date
	//gets json output from API and branch if an error occurred
	status, err := policyCurrent.req(url)
	if err != nil {
		return time.Time{}, status, err
	}
	//branch if output from API is empty and return error
	if policyCurrent.isEmpty() {
		err = errors.New("policyCurrent validation: policyCurrent is empty")
		return time.Time{}, http.StatusNotFound, err
	}
	return currentTime, http.StatusOK, nil
}
// req will request from API based on URL.
func (policyCurrent *PolicyCurrent) req(url string) (int, error) {
	//gets raw data from API and branch if an error occurred
	data, status, err := requestData(url)
	if err != nil {
		return status, err
	}
	//convert raw data to JSON and branch if an error occurred
	err = json.Unmarshal(data, &policyCurrent)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
// isEmpty checks if PolicyCurrent is empty.
func (policyCurrent *PolicyCurrent) isEmpty() bool {
    return policyCurrent.Stringencydata.Stringency == nil
}

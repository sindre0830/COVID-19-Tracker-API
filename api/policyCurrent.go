package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// PolicyCurrent stores data about current COVID policies based on a country.
//
// Functionality: Get, req, isEmpty
type PolicyCurrent struct {
	Policyactions []struct {
		PolicyTypeCode          string       `json:"policy_type_code"`
		PolicyTypeDisplay       string       `json:"policy_type_display"`
		Policyvalue             interface{}  `json:"policyvalue"`
		PolicyvalueActual       *interface{} `json:"policyvalue_actual"`
		Flagged                 interface{}  `json:"flagged"`
		IsGeneral               *interface{} `json:"is_general"`
		Notes                   interface{}  `json:"notes"`
		FlagValueDisplayField   string       `json:"flag_value_display_field"`
		PolicyValueDisplayField string       `json:"policy_value_display_field"`
	} `json:"policyActions"`
	Stringencydata struct {
		DateValue        *interface{} `json:"date_value"`
		CountryCode      *interface{} `json:"country_code"`
		Confirmed        *interface{} `json:"confirmed"`
		Deaths           *interface{} `json:"deaths"`
		StringencyActual *interface{} `json:"stringency_actual"`
		Stringency       *interface{} `json:"stringency"`
		Msg              *interface{} `json:"msg"`
	} `json:"stringencyData"`
}
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

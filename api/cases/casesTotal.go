package cases

import (
	"encoding/json"
	"errors"
	"main/api"
	"net/http"
)

// CasesTotal stores all data about COVID cases.
//
// Functionality: get, req, isEmpty
type CasesTotal struct {
	All struct {
		Confirmed         int          `json:"confirmed"`
		Recovered         int          `json:"recovered"`
		Deaths            int          `json:"deaths"`
		Country           string       `json:"country"`
		Population        int          `json:"population"`
		SqKmArea          int          `json:"sq_km_area"`
		LifeExpectancy    *interface{} `json:"life_expectancy"`
		ElevationInMeters int          `json:"elevation_in_meters"`
		Continent         string       `json:"continent"`
		Abbreviation      string       `json:"abbreviation"`
		Location          string       `json:"location"`
		Iso               int          `json:"iso"`
		CapitalCity       string       `json:"capital_city"`
		Lat               string       `json:"lat"`
		Long              string       `json:"long"`
		Updated           string       `json:"updated"`
	} `json:"all"`
}
// get will update CasesTotal based on input.
func (casesTotal *CasesTotal) Get(country string) (int, error) {
	url := "https://covid-api.mmediagroup.fr/v1/cases?country=" + country
	//gets json output from API and branch if an error occurred
	status, err := casesTotal.req(url)
	if err != nil {
		return status, err
	}
	//branch if output from API is empty and return error
	if casesTotal.isEmpty() {
		err = errors.New("casesTotal validation: casesTotal is empty")
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
// req will request from API based on URL.
func (casesTotal *CasesTotal) req(url string) (int, error) {
	//gets raw data from API and branch if an error occurred
	data, status, err := api.RequestData(url)
	if err != nil {
		return status, err
	}
	//convert raw data to JSON and branch if an error occurred
	err = json.Unmarshal(data, &casesTotal)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
// isEmpty checks if CasesTotal is empty.
func (casesTotal *CasesTotal) isEmpty() bool {
    return casesTotal.All.Country == ""
}

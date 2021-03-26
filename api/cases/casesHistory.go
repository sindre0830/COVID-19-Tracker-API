package cases

import (
	"encoding/json"
	"errors"
	"main/api"
	"net/http"
)

// CasesHistory structure stores historcial data about COVID cases for a country.
//
// Functionality: Get, req, isEmpty
type CasesHistory struct {
	All struct {
		Country           string         `json:"country"`
		Population        int            `json:"population"`
		SqKmArea          int            `json:"sq_km_area"`
		LifeExpectancy    *interface{}   `json:"life_expectancy"`
		ElevationInMeters interface{}    `json:"elevation_in_meters"`
		Continent         string         `json:"continent"`
		Abbreviation      string         `json:"abbreviation"`
		Location          string         `json:"location"`
		Iso               int            `json:"iso"`
		CapitalCity       string         `json:"capital_city"`
		Dates             map[string]int `json:"dates"`
	} `json:"All"`
}

// Get will get data for structure.
func (casesHistory *CasesHistory) Get(target string, country string) (int, error) {
	url := "https://covid-api.mmediagroup.fr/v1/history?country=" + country + "&status=" + target
	//gets json output from API and branch if an error occurred
	status, err := casesHistory.req(url)
	if err != nil {
		return status, err
	}
	//branch if output from API is empty and return error
	if casesHistory.isEmpty() {
		err = errors.New("casesHistory validation: casesHistory is empty")
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

// req will request data from API.
func (casesHistory *CasesHistory) req(url string) (int, error) {
	//gets raw data from API and branch if an error occurred
	output, status, err := api.RequestData(url)
	if err != nil {
		return status, err
	}
	//convert raw data to JSON and branch if an error occurred
	err = json.Unmarshal(output, &casesHistory)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// isEmpty checks if structure is empty.
func (casesHistory *CasesHistory) isEmpty() bool {
    return casesHistory.All.Country == ""
}

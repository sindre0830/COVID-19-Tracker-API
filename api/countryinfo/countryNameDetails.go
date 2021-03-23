package countryinfo

import (
	"encoding/json"
	"main/api"
	"net/http"
)

// CountryNameDetails structure stores information about country name.
//
// Functionality: Get, req
type CountryNameDetails []struct {
	Name       string `json:"name"`
	Alpha3Code string `json:"alpha3Code"`
}

// Get will get data for structure.
func (countryNameDetails *CountryNameDetails) Get(country string) (int, error) {
	url := "https://restcountries.eu/rest/v2/name/" + country + "?fields=name;alpha3Code"
	//gets json output from API and branch if an error occurred
	status, err := countryNameDetails.req(url)
	if err != nil {
		return status, err
	}
	return http.StatusOK, nil
}

// req will request data from API.
func (countryNameDetails *CountryNameDetails) req(url string) (int, error) {
	//gets raw data from API and branch if an error occurred
	data, status, err := api.RequestData(url)
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

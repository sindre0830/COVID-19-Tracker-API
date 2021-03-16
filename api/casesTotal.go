package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

// get will update object based on input.
func (object *casesTotal) get(country string) (int, error) {
	url := "https://covid-api.mmediagroup.fr/v1/cases?country=" + country
	//gets json output from API and branch if an error occurred
	status, err := object.req(url)
	if err != nil {
		return status, err
	}
	//branch if output from API is empty and return error
	if object.isEmpty() {
		err = errors.New("object validation: object is empty")
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
// req will request from API based on URL.
func (object *casesTotal) req(url string) (int, error) {
	//gets raw data from API and branch if an error occurred
	data, status, err := requestData(url)
	if err != nil {
		return status, err
	}
	//convert raw data to JSON and branch if an error occurred
	err = json.Unmarshal(data, &object)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
// isEmpty checks if object is empty.
func (object *casesTotal) isEmpty() bool {
    return object.All.Country == ""
}

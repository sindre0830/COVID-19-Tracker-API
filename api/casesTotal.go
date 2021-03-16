package api

import (
	"encoding/json"
	"errors"
)

func (object *casesTotal) get(country string) (int, error) {
	//url to API
	url := "https://covid-api.mmediagroup.fr/v1/cases?country=" + country
	status, err := object.req(url)
	if err != nil {
		return status, err
	}
	//branch if object is empty
	if object.isEmpty() {
		err = errors.New("object validation: object is empty")
		return 0, err
	}
	return 0, nil
}

func (object *casesTotal) req(url string) (int, error) {
	//gets raw output from API
	output, status, err := requestData(url)
	//branch if there is an error
	if err != nil {
		return status, err
	}
	//convert raw output to JSON
	err = json.Unmarshal(output, &object)
	return 0, err
}

func (object *casesTotal) isEmpty() bool {
    return object.All.Country == ""
}

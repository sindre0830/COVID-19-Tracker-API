package api

import (
	"encoding/json"
	"errors"
)

func (object *casesTotal) get(country string) error {
	//url to API
	url := "https://covid-api.mmediagroup.fr/v1/cases?country=" + country
	err := object.req(url)
	//branch if object is empty
	if object.isEmpty() {
		err = errors.New("object validation: object is empty, either country is mistyped or doesn't exist in our database")
		return err
	}
	return err
}

func (object *casesTotal) req(url string) error {
	//gets raw output from API
	output, err := requestData(url)
	//branch if there is an error
	if err != nil {
		return err
	}
	//convert raw output to JSON
	err = json.Unmarshal(output, &object)
	return err
}

func (object *casesTotal) isEmpty() bool {
    return object.All.Country == ""
}

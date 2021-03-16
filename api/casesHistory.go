package api

import (
	"encoding/json"
	"errors"
)

func (object *casesHistory) get(country string, startDate string, endDate string) (int, int, int, error) {
	//url to API with confirmed cases
	url := "https://covid-api.mmediagroup.fr/v1/history?country=" + country + "&status=Confirmed"
	status, err := object.req(url)
	//branch if there is an error
	if err != nil {
		return 0, 0, status, err
	}
	//branch if object is empty
	if object.isEmpty() {
		err = errors.New("object validation: object is empty")
		return 0, 0, 0, err
	}
	confirmed := object.addCases(startDate, endDate)
	//url to API with recovered cases
	url = "https://covid-api.mmediagroup.fr/v1/history?country=" + country + "&status=Recovered"
	status, err = object.req(url)
	//branch if there is an error
	if err != nil {
		return 0, 0, status, err
	}
	recovered := object.addCases(startDate, endDate)
	return confirmed, recovered, 0, nil
}

func (object *casesHistory) isEmpty() bool {
    return object.All.Country == ""
}

func (object *casesHistory) addCases(startDate string, endDate string) int {
	n:= object.All.Dates[endDate] - object.All.Dates[startDate]
	if n < 0 {
		n *= (-1)
	}
	return n
}

func (object *casesHistory) req(url string) (int, error) {
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

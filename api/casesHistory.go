package api

import (
	"encoding/json"
	"errors"
)

func (object *casesHistory) addCases(startDate string, endDate string) int {
	n:= object.All.Dates[endDate] - object.All.Dates[startDate]
	if n < 0 {
		n *= (-1)
	}
	return n
}

func (object *casesHistory) get(country string, startDate string, endDate string) (int, int, error) {
	//url to API with confirmed cases
	url := "https://covid-api.mmediagroup.fr/v1/history?country=" + country + "&status=Confirmed"
	err := object.req(url)
	//branch if object is empty
	if object.isEmpty() {
		err = errors.New("object validation: object is empty, either country is mistyped or doesn't exist in our database")
		return 0, 0, err
	}
	//branch if there is an error
	if err != nil {
		return 0, 0, err
	}
	confirmed := object.addCases(startDate, endDate)
	//url to API with recovered cases
	url = "https://covid-api.mmediagroup.fr/v1/history?country=" + country + "&status=Recovered"
	err = object.req(url)
	//branch if there is an error
	if err != nil {
		return 0, 0, err
	}
	recovered := object.addCases(startDate, endDate)
	return confirmed, recovered, nil
}

func (object *casesHistory) req(url string) error {
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

func (object *casesHistory) isEmpty() bool {
    return object.All.Country == ""
}

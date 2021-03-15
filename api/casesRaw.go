package api

import (
	"encoding/json"
)

type casesRaw struct {
	All struct {
		Confirmed         int    `json:"confirmed"`
		Recovered         int    `json:"recovered"`
		Deaths            int    `json:"deaths"`
		Country           string `json:"country"`
		Population        int    `json:"population"`
		SqKmArea          int    `json:"sq_km_area"`
		LifeExpectancy    string `json:"life_expectancy"`
		ElevationInMeters int    `json:"elevation_in_meters"`
		Continent         string `json:"continent"`
		Abbreviation      string `json:"abbreviation"`
		Location          string `json:"location"`
		Iso               int    `json:"iso"`
		CapitalCity       string `json:"capital_city"`
		Lat               string `json:"lat"`
		Long              string `json:"long"`
		Updated           string `json:"updated"`
	} `json:"all"`
}

func (object *casesRaw) get(country string) error {
	//url to API
	url := "https://covid-api.mmediagroup.fr/v1/cases?country=" + country
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

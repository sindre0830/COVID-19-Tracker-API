package api

import (
	"encoding/json"
)

type casesHistory struct {
	All struct {
		Country           string 		 `json:"country"`
		Population        int    		 `json:"population"`
		SqKmArea          int    		 `json:"sq_km_area"`
		LifeExpectancy    string 		 `json:"life_expectancy"`
		ElevationInMeters int    		 `json:"elevation_in_meters"`
		Continent         string 		 `json:"continent"`
		Abbreviation      string 		 `json:"abbreviation"`
		Location          string 		 `json:"location"`
		Iso               int    		 `json:"iso"`
		CapitalCity       string 		 `json:"capital_city"`
		Dates 		      map[string]int `json:"dates"`
	} `json:"All"`
}

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

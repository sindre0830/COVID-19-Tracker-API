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

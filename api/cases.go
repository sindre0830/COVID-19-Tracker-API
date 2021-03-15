package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Cases struct {
	Country              string  `json:"country"`
	Continent            string  `json:"continent"`
	Scope                string  `json:"scope"`
	Confirmed            int     `json:"confirmed"`
	Recovered            int     `json:"recovered"`
	PopulationPercentage float32 `json:"population_percentage"`
}

func (object *Cases) Handler(w http.ResponseWriter, r *http.Request) {
	err := object.get("Norway", "2020-12-01", "2021-01-31")
	//branch if there is an error
	if err != nil {
		fmt.Println(err)
	}
	//set header to JSON
	w.Header().Set("Content-Type", "application/json")
	//send output to user
	err = json.NewEncoder(w).Encode(object)
	//branch if there is an error
	if err != nil {
		fmt.Println(err)
	}
}

func (object *Cases) getTotal(country string) error {
	var data casesTotal
	err := data.get(country)
	//branch if there is an error
	if err != nil {
		return err
	}
	object.Country = data.All.Country
	object.Continent = data.All.Continent
	object.Scope = "total"
	object.Confirmed = data.All.Confirmed
	object.Recovered = data.All.Recovered
	object.PopulationPercentage = 0.00
	return nil
}

func (object *Cases) getHistory(country string, startDate string, endDate string) error {
	var data casesHistory
	confirmed, recovered, err := data.get(country, startDate, endDate)
	//branch if there is an error
	if err != nil {
		return err
	}
	object.Country = data.All.Country
	object.Continent = data.All.Continent
	object.Scope = startDate + "-" + endDate
	object.Confirmed = confirmed
	object.Recovered = recovered
	object.PopulationPercentage = 0.00
	return nil
}

func (object *Cases) get(country string, startDate string, endDate string) error {
	if startDate == "" {
		err := object.getTotal(country)
		//branch if there is an error
		if err != nil {
			return err
		}
	} else {
		err := object.getHistory(country, startDate, endDate)
		//branch if there is an error
		if err != nil {
			return err
		}
	}
	return nil
}

package api

import (
	"encoding/json"
	"main/debug"
	"main/fun"
	"net/http"
	"net/url"
	"strings"
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
	//split URL path by '/'
	arrURL := strings.Split(r.URL.Path, "/")
	//branch if there is an error
	if len(arrURL) != 5 {
		debug.UpdateErrorMessage(
			http.StatusBadRequest, 
			"Cases.Handler() -> Checking length of URL",
			"url validation: either too many or too few arguments in url path",
			"Path format. Expected format: '.../country?scope=start_at-end_at' ('?scope=start_at-end_at' is optional). Example: '.../norway?scope=2020-01-20-2021-02-01'.",
		)
		debug.PrintErrorInformation(w)
		return
	}
	country := fun.ConvertCountry(arrURL[4])
	err := fun.ValidateCountry(country)
	if err != nil {
		debug.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"Cases.Handler() -> ValidateCountry() -> Checking if inputted country is correct",
			err.Error(),
			"Country format. Expected format: '.../country...'. Example '.../Norway...'",
		)
		debug.PrintErrorInformation(w)
		return
	}
	//set default scope to nil (total)
	startDate := ""
	endDate := ""
	//get all parameters from URL
	arrPathParameters, err := url.ParseQuery(r.URL.RawQuery)
	//branch if there is an error
	if err != nil {
		debug.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"Cases.Handler() -> Getting URL field (...?scope=start_at-end_at)",
			err.Error(),
			"Unknown",
		)
		debug.PrintErrorInformation(w)
		return
	}
	//branch if any parameters exist
	if len(arrPathParameters) > 0 {
		//branch if field 'scope' exist
		if targetParameter, ok := arrPathParameters["scope"]; ok {
			dates := targetParameter[0]
			err := fun.ValidateDates(dates)
			//branch if there is an error
			if err != nil {
				debug.UpdateErrorMessage(
					http.StatusBadRequest, 
					"Cases.Handler() -> Checking if inputed dates are valid",
					err.Error(),
					"Date format. Expected format: '...?start_at-end_at' (YYYY-MM-DD-YYYY-MM-DD). Example: '...?2020-01-20-2021-02-01'",
				)
				debug.PrintErrorInformation(w)
				return
			}
			//set start- and end date variables
			startDate = dates[:10]
			endDate = dates[11:]
		//branch if there is an error
		} else {
			debug.UpdateErrorMessage(
				http.StatusBadRequest, 
				"Cases.Handler() -> Validating path parameters",
				"path validation: fields in URL used, but doesn't contain 'scope'",
				"Wrong field, or typo. Expected format: '...?scope=start_at-end_at'. Example: '...?scope=2020-01-20-2021-02-01'.",
			)
			debug.PrintErrorInformation(w)
			return
		}
	}
	//get data based on country and scope
	err = object.get(country, startDate, endDate)
	//branch if there is an error
	if err != nil {
		debug.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"Cases.Handler() -> Cases.get() -> Getting covid cases data",
			err.Error(),
			"Unknown",
		)
		debug.PrintErrorInformation(w)
		return
	}
	//set header to JSON
	w.Header().Set("Content-Type", "application/json")
	//send output to user
	err = json.NewEncoder(w).Encode(object)
	//branch if something went wrong with output
	if err != nil {
		debug.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"Cases.Handler() -> Sending data to user",
			err.Error(),
			"Unknown",
		)
		debug.PrintErrorInformation(w)
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

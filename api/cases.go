package api

import (
	"encoding/json"
	"main/debug"
	"main/fun"
	"math"
	"net/http"
	"net/url"
	"strings"
)

// Handler will handle http request.
func (cases *Cases) Handler(w http.ResponseWriter, r *http.Request) {
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
			http.StatusBadRequest, 
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
	status, err := cases.get(country, startDate, endDate)
	//branch if there is an error
	if err != nil {
		reason := "Unknown"
		if status == http.StatusBadRequest {
			reason = "Country format. Either country doesn't exist in our database or it's mistyped"
		}
		debug.UpdateErrorMessage(
			status, 
			"Cases.Handler() -> Cases.get() -> Getting covid cases data",
			err.Error(),
			reason,
		)
		debug.PrintErrorInformation(w)
		return
	}
	//set header to JSON
	w.Header().Set("Content-Type", "application/json")
	//send output to user
	err = json.NewEncoder(w).Encode(cases)
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
// get will update Cases based on input.
func (cases *Cases) get(country string, startDate string, endDate string) (int, error) {
	//branch if scope parameter is used
	if startDate == "" {
		//get all available data and branch if an error occurred
		status, err := cases.getTotal(country)
		if err != nil {
			return status, err
		}
	} else {
		//get data between two dates and branch if an error occurred
		status, err := cases.getHistory(country, startDate, endDate)
		if err != nil {
			return status, err
		}
	}
	return http.StatusOK, nil
}
// getTotal will get all available data.
func (cases *Cases) getTotal(country string) (int, error) {
	var data casesTotal
	//get total cases and branch if an error occurred
	status, err := data.get(country)
	if err != nil {
		return status, err
	}
	//set data in cases
	cases.update(data.All.Country, data.All.Continent, "total", data.All.Confirmed, data.All.Recovered, data.All.Population)
	return http.StatusOK, nil
}
// getHistory will get data between two dates.
func (cases *Cases) getHistory(country string, startDate string, endDate string) (int, error) {
	var data casesHistory
	//get cases between two dates and branch if an error occurred
	confirmed, recovered, status, err := data.get(country, startDate, endDate)
	if err != nil {
		return status, err
	}
	//set data in cases
	cases.update(data.All.Country, data.All.Continent, startDate + "-" + endDate, confirmed, recovered, data.All.Population)
	return http.StatusOK, nil
}
// update sets new data in cases.
func (cases *Cases) update(country string, continent string, scope string, confirmed int, recovered int, population int) {
	cases.Country = country
	cases.Continent = continent
	cases.Scope = scope
	cases.Confirmed = confirmed
	cases.Recovered = recovered
	//https://yourbasic.org/golang/round-float-2-decimal-places/#float-to-float
	cases.PopulationPercentage = (math.Round((float64(confirmed) / float64(population)) * 100) / 100)
}

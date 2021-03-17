package api

import (
	"encoding/json"
	"main/debug"
	"main/fun"
	"net/http"
	"time"
)

// Policy stores data about COVID policies based on user input.
//
// Functionality: Handler, get, getTotal, getHistory, update
type Policy struct {
	Country    string  `json:"country"`
	Scope      string  `json:"scope"`
	Stringency float64 `json:"stringency"`
	Trend	   float64 `json:"trend"`
	Update 	   string  `json:"update"`
}
// Handler will handle http request for COVID policies.
func (policy *Policy) Handler(w http.ResponseWriter, r *http.Request) {
	//parse url and branch if an error occurred
	country, scope, status, err := fun.ParseURL(r.URL)
	if err != nil {
		debug.UpdateErrorMessage(
			status, 
			"Policy.Handler() -> Parsing URL",
			err.Error(),
			"URL format. Expected format: '.../country?start_at-end_at' (YYYY-MM-DD-YYYY-MM-DD). Example: '.../norway?2020-01-20-2021-02-01'",
		)
		debug.PrintErrorInformation(w)
		return
	}
	//validate country name and branch if an error occurred
	err = fun.ValidateCountry(country)
	if err != nil {
		debug.UpdateErrorMessage(
			http.StatusBadRequest,
			"Cases.Handler() -> ValidatingCountry() -> Checking if inputed country is valid",
			err.Error(),
			"Country format. Expected format: '.../country'. Example: '.../norway'",
		)
		debug.PrintErrorInformation(w)
		return
	}
	//convert to required syntax
	country = fun.ConvertCountry(country)
	//set default start- and end date variables (total) and check if user inputted scope
	startDate := ""
	endDate := ""
	if len(scope) > 0 {
		//validate scope and branch if an error occurred
		err = fun.ValidateDates(scope)
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
		startDate = scope[:10]
		endDate = scope[11:]
	}
	//get data based on country and scope
	status, err = policy.get(country, startDate, endDate)
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
	err = json.NewEncoder(w).Encode(policy)
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
// get will update Policy based on input.
func (policy *Policy) get(country string, startDate string, endDate string) (int, error) {
	var countryCode CountryNameDetails
	countryCode.Get(country)
	//branch if scope parameter is used
	if startDate == "" {
		//get all available data and branch if an error occurred
		status, err := policy.getCurrent(countryCode[0].Alpha3Code)
		if err != nil {
			return status, err
		}
	} else {
		//get data between two dates and branch if an error occurred
		status, err := policy.getHistory(countryCode[0].Alpha3Code, startDate, endDate)
		if err != nil {
			return status, err
		}
	}
	policy.Country = country
	return http.StatusOK, nil
}
// getCurrent will get current available COVID policies.
func (policy *Policy) getCurrent(country string) (int, error) {
	var data PolicyCurrent
	//get total cases and branch if an error occurred
	updated, status, err := data.Get(country)
	if err != nil {
		return status, err
	}
	//set data in cases
	policy.update("total", data.Stringencydata.Stringency, 0, updated)
	return http.StatusOK, nil
}
// getHistory will get COVID policies between two dates.
func (policy *Policy) getHistory(country string, startDate string, endDate string) (int, error) {
	var data PolicyHistory
	//get total cases and branch if an error occurred
	trend, status, err := data.Get(country, startDate, endDate)
	if err != nil {
		return status, err
	}
	currentTime := time.Now()
	//set data in cases
	policy.update(startDate + "-" + endDate, data.Data[country][endDate].Stringency, trend, currentTime.String())
	return http.StatusOK, nil
}
// update sets new data in cases.
func (policy *Policy) update(scope string, stringency float64, trend float64, update string) {
	policy.Scope = scope
	policy.Stringency = stringency
	policy.Trend = trend
	policy.Update = update
}
